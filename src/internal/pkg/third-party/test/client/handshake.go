package watest

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"time"

	"google.golang.org/protobuf/proto"

	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/socket"
	"go.mau.fi/whatsmeow/util/keys"
)

// doHandshake implements the Noise_XX_25519_AESGCM_SHA256 handshake for the WhatsApp web API.
func doHandshake(fs *socket.FrameSocket, ephemeralKP keys.KeyPair, priv string, user string, index int) error {
	nh := socket.NewNoiseHandshake()
	nh.Start(socket.NoiseStartPattern, fs.Header)
	nh.Authenticate(ephemeralKP.Pub[:])
	data, err := proto.Marshal(&waProto.HandshakeMessage{
		ClientHello: &waProto.HandshakeClientHello{
			Ephemeral: ephemeralKP.Pub[:],
		},
	})
	if err != nil {
		return fmt.Errorf("failed to marshal handshake message: %w", err)
	}
	err = fs.SendFrame(data)
	if err != nil {
		return fmt.Errorf("failed to send handshake message: %w", err)
	}
	var resp []byte
	select {
	case resp = <-fs.Frames:
	case <-time.After(20 * time.Second):
		return fmt.Errorf("timed out waiting for handshake response")
	}
	var handshakeResponse waProto.HandshakeMessage
	err = proto.Unmarshal(resp, &handshakeResponse)
	if err != nil {
		return fmt.Errorf("failed to unmarshal handshake response: %w", err)
	}
	serverEphemeral := handshakeResponse.GetServerHello().GetEphemeral()
	serverStaticCiphertext := handshakeResponse.GetServerHello().GetStatic()
	certificateCiphertext := handshakeResponse.GetServerHello().GetPayload()
	if len(serverEphemeral) != 32 || serverStaticCiphertext == nil || certificateCiphertext == nil {
		return fmt.Errorf("missing parts of handshake response")
	}
	serverEphemeralArr := *(*[32]byte)(serverEphemeral)

	nh.Authenticate(serverEphemeral)
	err = nh.MixSharedSecretIntoKey(*ephemeralKP.Priv, serverEphemeralArr)
	if err != nil {
		return fmt.Errorf("failed to mix server ephemeral key in: %w", err)
	}

	staticDecrypted, err := nh.Decrypt(serverStaticCiphertext)
	if err != nil {
		return fmt.Errorf("failed to decrypt server static ciphertext: %w", err)
	} else if len(staticDecrypted) != 32 {
		return fmt.Errorf("unexpected length of server static plaintext %d (expected 32)", len(staticDecrypted))
	}
	err = nh.MixSharedSecretIntoKey(*ephemeralKP.Priv, *(*[32]byte)(staticDecrypted))
	if err != nil {
		return fmt.Errorf("failed to mix server static key in: %w", err)
	}

	certDecrypted, err := nh.Decrypt(certificateCiphertext)
	if err != nil {
		return fmt.Errorf("failed to decrypt noise certificate ciphertext: %w", err)
	}
	var cert waProto.NoiseCertificate
	err = proto.Unmarshal(certDecrypted, &cert)
	if err != nil {
		return fmt.Errorf("failed to unmarshal noise certificate: %w", err)
	}
	certDetailsRaw := cert.GetDetails()
	certSignature := cert.GetSignature()
	if certDetailsRaw == nil || certSignature == nil {
		return fmt.Errorf("missing parts of noise certificate")
	}
	var certDetails waProto.NoiseCertificate_Details
	err = proto.Unmarshal(certDetailsRaw, &certDetails)
	if err != nil {
		return fmt.Errorf("failed to unmarshal noise certificate details: %w", err)
	} else if !bytes.Equal(certDetails.GetKey(), staticDecrypted) {
		return fmt.Errorf("cert key doesn't match decrypted static")
	}

	//

	//dbLog := waLog.Stdout("Database", logLevel, true)
	//storeContainer, err := sqlstore.New("sqlite3",
	//	fmt.Sprintf("file:%s?_foreign_keys=off", "/root/wa/go-socialapp/src/storages/whatsapp.db"), dbLog)
	////storeContainer, err := sqlstore.New(*dbDialect, *dbAddress, dbLog)
	//if err != nil {
	//	log.Errorf("Failed to connect to database: %v", err)
	//	return err
	//}
	//device, err := storeContainer.GetFirstDevice()
	//if err != nil {
	//	log.Errorf("Failed to get device: %v", err)
	//	return err
	//}
	//var noisePriv, identityPriv, preKeyPriv, preKeySig []byte
	//keys.NewKeyPairFromPrivateKey(*(*[32]byte)(noisePriv))

	de, err := base64.StdEncoding.DecodeString(priv)

	spew.Dump("--err=", err)
	spew.Dump("--priv=", priv)

	noisekey := keys.NewKeyPairFromPrivateKey(*(*[32]byte)(de))

	encryptedPubkey := nh.Encrypt(noisekey.Pub[:])
	err = nh.MixSharedSecretIntoKey(*noisekey.Priv, serverEphemeralArr)
	if err != nil {
		return fmt.Errorf("failed to mix noise private key in: %w", err)
	}

	clientFinishPayloadBytes, err := proto.Marshal(getLoginPayload(user, index))
	if err != nil {
		return fmt.Errorf("failed to marshal client finish payload: %w", err)
	}
	encryptedClientFinishPayload := nh.Encrypt(clientFinishPayloadBytes)
	data, err = proto.Marshal(&waProto.HandshakeMessage{
		ClientFinish: &waProto.HandshakeClientFinish{
			Static:  encryptedPubkey,
			Payload: encryptedClientFinishPayload,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to marshal handshake finish message: %w", err)
	}
	err = fs.SendFrame(data)
	if err != nil {
		return fmt.Errorf("failed to send handshake finish message: %w", err)
	}

	ns, err := nh.Finish(fs, handleFrame, onDisconnect)
	if err != nil {
		return fmt.Errorf("failed to create noise socket: %w", err)
	}

	spew.Dump("---ns--,", ns.IsConnected())
	//cli.socket = ns

	return nil
}
