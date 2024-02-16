package fixtures

import (
	"github.com/maguro-alternative/remake_bot/pkg/db"

	"testing"
)

type Fixture struct {
	PermissionsCodes   []*PermissionsCode
	PermissionsIDs     []*PermissionsID
	LineChannels       []*LineChannel
	LineNgTypes        []*LineNgType
	LineNgDiscordIDs   []*LineNgDiscordID
	LineBots           []*LineBot
	LineBotIvs         []*LineBotIv
	VcSignalChannels   []*VcSignalChannel
	VcSignalNgIDs      []*VcSignalNgID
	VcSignalMentionIDs []*VcSignalMentionID
	Webhooks           []*Webhook
	WebhookMentions    []*WebhookMention
	WebhookWords       []*WebhookWord

	DBv1 db.Driver
}

func (f *Fixture) Build(t *testing.T, modelConnectors ...*ModelConnector) *Fixture {

	for _, modelConnector := range modelConnectors {
		modelConnector.addToFixtureAndConnect(t, f)
	}

	return f
}

type ModelConnector struct {
	Model interface{}

	// 定義されるべきコールバック
	setter       func()
	addToFixture func(t *testing.T, f *Fixture)
	connect      func(t *testing.T, f *Fixture, connectingModel interface{})
	insertTable  func(t *testing.T, f *Fixture)

	// 状態
	addedToFixture bool
	connectings    []*ModelConnector
}

func (mc *ModelConnector) Connect(connectors ...*ModelConnector) *ModelConnector {
	mc.connectings = append(mc.connectings, connectors...)
	return mc // メソッドチェーンで記述できるようにする
}

func (mc *ModelConnector) addToFixtureAndConnect(t *testing.T, fixture *Fixture) {
	if mc.addedToFixture {
		return
	}

	if mc.addToFixture == nil {
		// addToFixtureは必ずセットされている必要がある
		t.Fatalf("addToFixture field of %T is not properly initialized", mc.Model)
	}
	// このモデルをfixtureに追加する
	mc.setter()
	mc.insertTable(t, fixture)
	mc.addToFixture(t, fixture)

	for _, modelConnector := range mc.connectings {
		if mc.connect == nil {
			// どのモデルとも接続できない場合はconnectをnilにできる
			t.Fatalf("%T cannot be connected to %T", modelConnector.Model, mc.Model)
		}

		mc.connect(t, fixture, modelConnector.Model)

		modelConnector.addToFixtureAndConnect(t, fixture)
	}

	mc.addedToFixture = true
}
