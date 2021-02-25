package discordgo

import (
	"encoding/json"
	"time"
)

type GatewayPayload struct {
	Op int             `json:"op"`
	D  json.RawMessage `json:"d,omitempty"`
	S  int             `json:"s,omitempty"`
	T  string          `json:"t,omitempty"`
}

type IdentityPayload struct {
	Token      string                       `json:"token"`
	Properties IdentityConnectionProperties `json:"properties"`
	Intents    int                          `json:"intents"`
}

type IdentityConnectionProperties struct {
	Os      string `json:"$os"`
	Browser string `json:"$browser"`
	Device  string `json:"$device"`
}

type HelloPayload struct {
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
}

type ResumePayload struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Sequence  int    `json:"seq"`
}

type Message struct {
	Id                string              `json:"id,omitempty"`
	ChannelID         string              `json:"channel_id,omitempty"`
	GuildID           string              `json:"guild_id,omitempty"`
	Author            *User               `json:"author,omitempty"`
	Content           string              `json:"content,omitempty"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
	EditedTimestamp   *time.Time           `json:"edited_timestamp,omitempty"`
	Tts               bool                `json:"tts,omitempty"`
	MentionsEveryone  bool                `json:"mention_everyone,omitempty"`
	UserMentions      []*User             `json:"mentions,omitempty"`
	RoleMentions      []string            `json:"mention_roles,omitempty"`
	ChannelMentions   []*ChannelMention   `json:"mention_channels,omitempty"`
	Attachments       []*Attachment       `json:"attachments,omitempty"`
	Embeds            []*Embed            `json:"embeds,omitempty"`
	Reactions         []*Reaction         `json:"reactions,omitempty"`
	Pinned            bool                `json:"pinned,omitempty"`
	WebhookID         string              `json:"webhook_id,omitempty"`
	Type              int                 `json:"type,omitempty"`
	Activity          *MessageActivity    `json:"activity,omitempty"`
	Application       *MessageApplication `json:"application,omitempty"`
	MessageReference  *MessageReference   `json:"message_reference,omitempty"`
	Flags             int                 `json:"flags,omitempty"`
	ReferencedMessage *Message            `json:"referenced_message,omitempty"`
}

type User struct {
	Id            string `json:"id,omitempty"`
	Username      string `json:"username,omitempty"`
	Discriminator string `json:"discriminator,omitempty"`
	Avatar        string `json:"avatar,omitempty"`
	Bot           bool   `json:"bot,omitempty"`
	System        bool   `json:"system,omitempty"`
	MfaEnabled    bool   `json:"mfa_enabled,omitempty"`
	Locale        string `json:"locale,omitempty"`
	Verified      bool `json:"verified,omitempty"`
	Email         string `json:"email,omitempty"`
	Flags         int    `json:"flags,omitempty"`
	PremiumType   int    `json:"premium_type,omitempty"`
	PublicFlags   int    `json:"public_flags,omitempty"`
}

type ChannelMention struct {
	Id      string `json:"id"`
	Type    int    `json:"type"`
	GuildID string `json:"guild_id"`
	Name    string `json:"name"`
}

type Attachment struct {
	Id       string `json:"id"`
	Filename string `json:"filename"`
	Size     int    `json:"size"`
	Url      string `json:"url"`
	ProxyUrl string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

type Embed struct {
	Title       string         `json:"title"`
	Type        string         `json:"type"`
	Description string         `json:"description"`
	Url         string         `json:"url"`
	Timestamp   string         `json:"timestamp"`
	Color       int            `json:"color"`
	Footer      *EmbedFooter    `json:"footer"`
	Image       *EmbedImage     `json:"image"`
	Thumbnail   *EmbedThumbnail `json:"thumbnail"`
	Video       *EmbedVideo     `json:"video"`
	Provider    *EmbedProvider  `json:"provider"`
	Author      *EmbedAuthor    `json:"author"`
	Fields      *[]EmbedField   `json:"fields"`
}

type EmbedFooter struct {
	Text         string `json:"text"`
	IconUrl      string `json:"icon_url"`
	ProxyIconUrl string `json:"proxy_icon_url"`
}

type EmbedImage struct {
	Url      string `json:"url"`
	ProxyUrl string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

type EmbedThumbnail struct {
	Url      string `json:"url"`
	ProxyUrl string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

type EmbedVideo struct {
	Url      string `json:"url"`
	ProxyUrl string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

type EmbedProvider struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type EmbedAuthor struct {
	Name         string `json:"name"`
	Url          string `json:"url"`
	IconUrl      string `json:"icon_url"`
	ProxyIconUrl string `json:"proxy_icon_url"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type Reaction struct {
	Count int   `json:"count"`
	Me    bool  `json:"me"`
	Emoji Emoji `json:"emoji"`
}

type Emoji struct {
	Id            string   `json:"id"`
	Name          string   `json:"name"`
	Roles         []string `json:"roles"`
	User          User     `json:"user"`
	RequireColons bool     `json:"require_colons"`
	Managed       bool     `json:"managed"`
	Animated      bool     `json:"animated"`
	Available     bool     `json:"available"`
}

type MessageActivity struct {
	Type    int    `json:"type"`
	PartyId string `json:"party_id"`
}

type MessageApplication struct {
	Id          string `json:"id"`
	CoverImage  string `json:"cover_image"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Name        string `json:"name"`
}

type MessageReference struct {
	MessageId string `json:"message_id,omitempty"`
	ChannelId string `json:"channel_id,omitempty"`
	GuildId   string `json:"guild_id,omitempty"`
}

type GetGatewayResponse struct {
	Url               string            `json:"url"`
	Shards            int               `json:"shards"`
	SessionStartLimit SessionStartLimit `json:"session_start_limit"`
}

type SessionStartLimit struct {
	Total          int `json:"total"`
	Remaining      int `json:"remaining"`
	ResetAfter     int `json:"reset_after"`
	MaxConcurrency int `json:"max_concurrency"`
}

type MessageDeleteEvent struct{
	Id string `json:"id"`
	ChannelID string `json:"channel_id"`
	GuildID string `json:"guild_id"`
}

type MessageDeleteBulkEvent struct{
	Ids []string `json:"ids"`
	ChannelID string `json:"channel_id"`
	GuildID string `json:"guild_id"`
}

type ReadyEvent struct{
	Version int `json:"v"`
	User User `json:"user"`
	SessionID string `json:"session_id"`
}

type Channel struct{
	Id string `json:"id,omitempty"`
	Type int `json:"type,omitempty"`
	GuildID string `json:"guild_id"`
	Overwrites []*ChannelOverwrite `json:"permission_overwrites,omitempty"`
	Position int `json:"position,omitempty"`
	Name string `json:"name,omitempty"`
	Topic string `json:"topic,omitempty"`
	Nsfw bool `json:"nsfw,omitempty"`
	LastMessageID string `json:"last_message_id,omitempty"`
	Bitrate int `json:"bitrate,omitempty"`
	UserLimit int `json:"user_limit,omitempty"`
	RateLimitPerUser int `json:"rate_limit_per_user,omitempty"`
	Recipients []*User `json:"recipients,omitempty"`
	Icon string `json:"icon,omitempty"`
	OwnerID string `json:"owner_id,omitempty"`
	ApplicationID string `json:"application_id,omitempty"`
	ParentID string `json:"parent_id,omitempty"`
	LastPinTimestamp time.Time `json:"last_pin_timestamp,omitempty"`
}

type ChannelOverwrite struct{
	Id string `json:"id,omitempty"`
	Type int `json:"type,omitempty"`
	Allow string `json:"allow,omitempty"`
	Deny string `json:"deny,omitempty"`
}