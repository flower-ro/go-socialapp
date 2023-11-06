package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/marmotedu/errors"
	"github.com/marmotedu/iam/pkg/log"
	"go-socialapp/internal/pkg/code"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/appstate"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
	"mime"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync/atomic"
	"time"
)

func handler(rawEvt interface{}) {
	switch evt := rawEvt.(type) {
	case *events.AppStateSyncComplete:
		if len(cli.Store.PushName) > 0 && evt.Name == appstate.WAPatchCriticalBlock {
			err := cli.SendPresence(types.PresenceAvailable)
			if err != nil {
				log.Warnf("Failed to send available presence: %v", err)
			} else {
				log.Infof("Marked self as available")
			}
		}
	case *events.PairSuccess:
		log.Infof("Successfully pair with %s", evt.ID.String())
		spew.Dump(evt)
		Broadcast <- BroadcastMessage{
			Code: "LOGIN_SUCCESS",
			//Message: fmt.Sprintf("Successfully pair with %s", evt.ID.String()),
			Result: evt.ID.String(),
		}
	case *events.LoggedOut:
		Broadcast <- BroadcastMessage{
			Code:   "LIST_DEVICES",
			Result: nil,
		}
	case *events.Connected, *events.PushNameSetting:
		if len(cli.Store.PushName) == 0 {
			return
		}

		// Send presence available when connecting and when the pushname is changed.
		// This makes sure that outgoing messages always have the right pushname.
		err := cli.SendPresence(types.PresenceAvailable)
		if err != nil {
			log.Warnf("Failed to send available presence: %v", err)
		} else {
			log.Infof("Marked self as available")
		}
	case *events.StreamReplaced:
		os.Exit(0)
	case *events.Message:
		metaParts := []string{fmt.Sprintf("pushname: %s", evt.Info.PushName), fmt.Sprintf("timestamp: %s", evt.Info.Timestamp)}
		if evt.Info.Type != "" {
			metaParts = append(metaParts, fmt.Sprintf("type: %s", evt.Info.Type))
		}
		if evt.Info.Category != "" {
			metaParts = append(metaParts, fmt.Sprintf("category: %s", evt.Info.Category))
		}
		if evt.IsViewOnce {
			metaParts = append(metaParts, "view once")
		}
		if evt.IsViewOnce {
			metaParts = append(metaParts, "ephemeral")
		}

		log.Infof("Received message %s from %s (%s): %+v", evt.Info.ID, evt.Info.SourceString(), strings.Join(metaParts, ", "), evt.Message)

		img := evt.Message.GetImageMessage()
		if img != nil {
			path, err := ExtractMedia(PathStorages, img)
			if err != nil {
				log.Errorf("Failed to download image: %v", err)
			} else {
				log.Infof("Image downloaded to %s", path)
			}
		}

		if WhatsappAutoReplyMessage != "" &&
			!isGroupJid(evt.Info.Chat.String()) &&
			!strings.Contains(evt.Info.SourceString(), "broadcast") {
			_, _ = cli.SendMessage(context.Background(), evt.Info.Sender, &waProto.Message{Conversation: proto.String(WhatsappAutoReplyMessage)})
		}

		if WhatsappWebhook != "" &&
			!strings.Contains(evt.Info.SourceString(), "broadcast") &&
			!isFromMySelf(evt.Info.SourceString()) {
			if err := forwardToWebhook(evt); err != nil {
				log.Errorf("Failed forward to webhook,err: %v", err)
			}
		}
	case *events.Receipt:
		if evt.Type == events.ReceiptTypeRead || evt.Type == events.ReceiptTypeReadSelf {
			log.Infof("%v was read by %s at %s", evt.MessageIDs, evt.SourceString(), evt.Timestamp)
		} else if evt.Type == events.ReceiptTypeDelivered {
			log.Infof("%s was delivered to %s at %s", evt.MessageIDs[0], evt.SourceString(), evt.Timestamp)
		}
	case *events.Presence:
		if evt.Unavailable {
			if evt.LastSeen.IsZero() {
				log.Infof("%s is now offline", evt.From)
			} else {
				log.Infof("%s is now offline (last seen: %s)", evt.From, evt.LastSeen)
			}
		} else {
			log.Infof("%s is now online", evt.From)
		}
	case *events.HistorySync:
		id := atomic.AddInt32(&historySyncID, 1)
		fileName := fmt.Sprintf("%s/history-%d-%s-%d-%s.json", PathStorages, startupTime, cli.Store.ID.String(), id, evt.Data.SyncType.String())
		file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.Errorf("Failed to open file to write history sync: %v", err)
			return
		}
		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		err = enc.Encode(evt.Data)
		if err != nil {
			log.Errorf("Failed to write history sync: %v", err)
			return
		}
		log.Infof("Wrote history sync to %s", fileName)
		_ = file.Close()
	case *events.AppState:
		log.Debugf("App state event: %+v / %+v", evt.Index, evt.SyncActionValue)
	}
}

func ExtractMedia(storageLocation string, mediaFile whatsmeow.DownloadableMessage) (extractedMedia ExtractedMedia, err error) {
	if mediaFile == nil {
		log.Info("Skip download because data is nil")
		return extractedMedia, nil
	}

	data, err := cli.Download(mediaFile)
	if err != nil {
		return extractedMedia, err
	}

	switch media := mediaFile.(type) {
	case *waProto.ImageMessage:
		extractedMedia.MimeType = media.GetMimetype()
		extractedMedia.Caption = media.GetCaption()
	case *waProto.AudioMessage:
		extractedMedia.MimeType = media.GetMimetype()
	case *waProto.VideoMessage:
		extractedMedia.MimeType = media.GetMimetype()
		extractedMedia.Caption = media.GetCaption()
	case *waProto.StickerMessage:
		extractedMedia.MimeType = media.GetMimetype()
	case *waProto.DocumentMessage:
		extractedMedia.MimeType = media.GetMimetype()
		extractedMedia.Caption = media.GetCaption()
	}

	extensions, _ := mime.ExtensionsByType(extractedMedia.MimeType)
	extractedMedia.MediaPath = fmt.Sprintf("%s/%d-%s%s", storageLocation, time.Now().Unix(), uuid.NewString(), extensions[0])
	err = os.WriteFile(extractedMedia.MediaPath, data, 0600)
	if err != nil {
		return extractedMedia, err
	}
	return extractedMedia, nil
}

// isFromMySelf is a helper function to check if the message is from my self (logged in account)
func isFromMySelf(jid string) bool {
	return extractPhoneNumber(jid) == extractPhoneNumber(cli.Store.ID.String())
}

// extractPhoneNumber is a helper function to extract the phone number from a JID
func extractPhoneNumber(jid string) string {
	regex := regexp.MustCompile(`\d+`)
	// Find all matches of the pattern in the JID
	matches := regex.FindAllString(jid, -1)
	// The first match should be the phone number
	if len(matches) > 0 {
		return matches[0]
	}
	// If no matches are found, return an empty string
	return ""
}

// isGroupJid is a helper function to check if the message is from group
func isGroupJid(jid string) bool {
	return strings.Contains(jid, "@g.us")
}

// forwardToWebhook is a helper function to forward event to webhook url
func forwardToWebhook(evt *events.Message) error {
	log.Infof("Forwarding event to webhook: %s", WhatsappWebhook)
	client := &http.Client{Timeout: 10 * time.Second}
	imageMedia := evt.Message.GetImageMessage()
	stickerMedia := evt.Message.GetStickerMessage()
	videoMedia := evt.Message.GetVideoMessage()
	audioMedia := evt.Message.GetAudioMessage()
	documentMedia := evt.Message.GetDocumentMessage()

	var message evtMessage
	message.Text = evt.Message.GetConversation()
	message.ID = evt.Info.ID
	if extendedMessage := evt.Message.ExtendedTextMessage.GetText(); extendedMessage != "" {
		message.Text = extendedMessage
	}

	var quotedmessage any
	if evt.Message.ExtendedTextMessage != nil {
		if conversation := evt.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.GetConversation(); conversation != "" {
			quotedmessage = conversation
		}
	}

	var forwarded any
	if evt.Message.ExtendedTextMessage != nil && evt.Message.ExtendedTextMessage.ContextInfo != nil {
		if isForwarded := evt.Message.ExtendedTextMessage.ContextInfo.GetIsForwarded(); !isForwarded {
			forwarded = nil
		}
	}

	var waReaction *evtReaction
	if evt.Message.ReactionMessage != nil {
		waReaction.Message = evt.Message.ReactionMessage.GetText()
		waReaction.ID = evt.Message.ReactionMessage.GetKey().GetId()
	}

	body := map[string]interface{}{
		"audio":          audioMedia,
		"contact":        evt.Message.GetContactMessage(),
		"document":       documentMedia,
		"forwarded":      forwarded,
		"from":           evt.Info.SourceString(),
		"image":          imageMedia,
		"list":           evt.Message.GetListMessage(),
		"live_location":  evt.Message.GetLiveLocationMessage(),
		"location":       evt.Message.GetLocationMessage(),
		"message":        message,
		"order":          evt.Message.GetOrderMessage(),
		"pushname":       evt.Info.PushName,
		"quoted_message": quotedmessage,
		"reaction":       waReaction,
		"sticker":        stickerMedia,
		"video":          videoMedia,
		"view_once":      evt.Message.GetViewOnceMessage(),
	}

	if imageMedia != nil {
		path, err := ExtractMedia(PathMedia, imageMedia)
		if err != nil {
			return errors.WithCode(code.ErrWebhook, fmt.Sprintf("Failed to download image: %v", err))
		}
		body["image"] = path
	}
	if stickerMedia != nil {
		path, err := ExtractMedia(PathMedia, stickerMedia)
		if err != nil {
			return errors.WithCode(code.ErrWebhook, fmt.Sprintf("Failed to download sticker: %v", err))
		}
		body["sticker"] = path
	}
	if videoMedia != nil {
		path, err := ExtractMedia(PathMedia, videoMedia)
		if err != nil {
			return errors.WithCode(code.ErrWebhook, fmt.Sprintf("Failed to download video: %v", err))
		}
		body["video"] = path
	}
	if audioMedia != nil {
		path, err := ExtractMedia(PathMedia, audioMedia)
		if err != nil {
			return errors.WithCode(code.ErrWebhook, fmt.Sprintf("Failed to download audio: %v", err))
		}
		body["audio"] = path
	}
	if documentMedia != nil {
		path, err := ExtractMedia(PathMedia, documentMedia)
		if err != nil {
			return errors.WithCode(code.ErrWebhook, fmt.Sprintf("Failed to download document: %v", err))
		}
		body["document"] = path
	}

	postBody, err := json.Marshal(body)
	if err != nil {
		return errors.WithCode(code.ErrWebhook, fmt.Sprintf("Failed to marshal body: %v", err))
	}

	req, err := http.NewRequest(http.MethodPost, WhatsappWebhook, bytes.NewBuffer(postBody))
	if err != nil {
		return errors.WithCode(code.ErrWebhook, fmt.Sprintf("error when create http object %v", err))
	}
	req.Header.Set("Content-Type", "application/json")
	if _, err = client.Do(req); err != nil {
		return errors.WithCode(code.ErrWebhook, fmt.Sprintf("error when submit webhook %v", err))
	}
	return nil
}
