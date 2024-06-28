const { JSDOM } = require('jsdom');

const dom = new JSDOM(`
    <!doctype html>
    <html>
        <body>
            <form id="form">
                <select name="newWebhookType1" id="newWebhookType1">
                    <option value="Type1">Type1</option>
                </select>
                <select name="newSubscriptionName1" id="newSubscriptionName1">
                    <option value="Subscription1">Subscription1</option>
                </select>
                <select name="newMemberMention1[]" id="newMemberMention1[]">
                    <option value="Member1">Member1</option>
                </select>
                <select name="newRoleMention1[]" id="newRoleMention1[]">
                    <option value="Role1">Role1</option>
                </select>
            </form>
        </body>
    </html>
`);
global.document = dom.window.document;
global.window = dom.window;

global.fetch = jest.fn();
container = global.document.getElementById('form');

const localStorageMock = (() => {
    let store = {};

    return {
        getItem(key) {
            return store[key] || null;
        },
        setItem(key, value) {
            store[key] = value.toString();
        },
        removeItem(key) {
            delete store[key];
        },
        clear() {
            store = {};
        },
    };
})();

Object.defineProperty(window, 'sessionStorage', {
    value: localStorageMock,
});

const { createJsonData } = require('./webhook'); // 関数をインポートするパスを適宜調整してください

describe('createJsonData function', () => {
    it('should correctly process form data into JSON', async () => {
        // モックデータの準備
        const formElements = [
            { name: 'newWebhookType1', value: 'Type1' },
            { name: 'newSubscriptionType1', value: 'Subscription1' },
            { name: 'newSubscriptionId1', value: 'SubId1' },
            // 他のフィールドも同様に追加
        ];
        const formData = new FormData();
        formElements.forEach(el => formData.append(el.name, el.value));

        // モックのFormDataを模倣するための関数
        formData.get = (key) => {
            return formElements.find(el => el.name === key)?.value || null;
        };
        formData.getAll = (key) => {
            return formElements.filter(el => el.name === key).map(el => el.value);
        };

        // 関数の実行
        const result = await createJsonData(formElements, formData);

        // 期待値との比較
        expect(result).toBeDefined();
        expect(JSON.parse(result)).toEqual({
            newWebhooks: [
                {
                    webhookId: 'Type1',
                    subscriptionType: 'Subscription1',
                    subscriptionId: 'SubId1',
                    mentionRoles: [],
                    mentionUsers: [],
                    ngOrWords: [],
                    ngAndWords: [],
                    searchOrWords: [],
                    searchAndWords: [],
                    mentionOrWords: [],
                    mentionAndWords: [],
                }
            ],
            updateWebhooks: [] // updateWebhooksの期待値も同様に設定
        });
    });

    it('should correctly process form data into JSON with multiple mentions', async () => {
        // モックデータの準備
        const formElements = [
            { name: 'newWebhookType1', value: 'Type1' },
            { name: 'newSubscriptionType1', value: 'Subscription1' },
            { name: 'newSubscriptionId1', value: 'SubId1' },
            { name: 'newMemberMention1[]', value: 'Member1' },
            { name: 'newMemberMention1[]', value: 'Member2' },
            { name: 'newRoleMention1[]', value: 'Role1' },
            { name: 'newRoleMention1[]', value: 'Role2' },
            // 他のフィールドも同様に追加
        ];
        const formData = new FormData();
        formElements.forEach(el => formData.append(el.name, el.value));

        // モックのFormDataを模倣するための関数
        formData.get = (key) => {
            return formElements.find(el => el.name === key)?.value || null;
        };
        formData.getAll = (key) => {
            return formElements.filter(el => el.name === key).map(el => el.value);
        };

        // 関数の実行
        const result = await createJsonData(formElements, formData);

        // 期待値との比較
        expect(result).toBeDefined();
        expect(JSON.parse(result)).toEqual({
            newWebhooks: [
                {
                    webhookId: 'Type1',
                    subscriptionType: 'Subscription1',
                    subscriptionId: 'SubId1',
                    mentionRoles: ['Role1', 'Role2'],
                    mentionUsers: ['Member1', 'Member2'],
                    ngOrWords: [],
                    ngAndWords: [],
                    searchOrWords: [],
                    searchAndWords: [],
                    mentionOrWords: [],
                    mentionAndWords: [],
                }
            ],
            updateWebhooks: [] // updateWebhooksの期待値も同様に設定
        });
    });

    it('should correctly process form data into JSON with multiple mentions and NG words', async () => {
        // モックデータの準備
        const formElements = [
            { name: 'newWebhookType1', value: 'Type1' },
            { name: 'newSubscriptionType1', value: 'Subscription1' },
            { name: 'newSubscriptionId1', value: 'SubId1' },
            { name: 'newMemberMention1[]', value: 'Member1' },
            { name: 'newMemberMention1[]', value: 'Member2' },
            { name: 'newRoleMention1[]', value: 'Role1' },
            { name: 'newRoleMention1[]', value: 'Role2' },
            { name: 'newNgOrWord1[]', value: 'NGWord1' },
            { name: 'newNgOrWord1[]', value: 'NGWord2' },
            { name: 'newNgAndWord1[]', value: 'NGWord3' },
            { name: 'newNgAndWord1[]', value: 'NGWord4' },
            // 他のフィールドも同様に追加
        ];
        const formData = new FormData();
        formElements.forEach(el => formData.append(el.name, el.value));

        // モックのFormDataを模倣するための関数
        formData.get = (key) => {
            return formElements.find(el => el.name === key)?.value || null;
        };
        formData.getAll = (key) => {
            return formElements.filter(el => el.name === key).map(el => el.value);
        };

        // 関数の実行
        const result = await createJsonData(formElements, formData);

        // 期待値との比較
        expect(result).toBeDefined();
        expect(JSON.parse(result)).toEqual({
            newWebhooks: [
                {
                    webhookId: 'Type1',
                    subscriptionType: 'Subscription1',
                    subscriptionId: 'SubId1',
                    mentionRoles: ['Role1', 'Role2'],
                    mentionUsers: ['Member1', 'Member2'],
                    ngOrWords: ['NGWord1', 'NGWord2'],
                    ngAndWords: ['NGWord3', 'NGWord4'],
                    searchOrWords: [],
                    searchAndWords: [],
                    mentionOrWords: [],
                    mentionAndWords: [],
                }
            ],
            updateWebhooks: [] // updateWebhooksの期待値も同様に設定
        });
    });

    it('should correctly process form data into JSON with multiple mentions, NG words, and search words', async () => {
        // モックデータの準備
        const formElements = [
            { name: 'newWebhookType1', value: 'Type1' },
            { name: 'newSubscriptionType1', value: 'Subscription1' },
            { name: 'newSubscriptionId1', value: 'SubId1' },
            { name: 'newMemberMention1[]', value: 'Member1' },
            { name: 'newMemberMention1[]', value: 'Member2' },
            { name: 'newRoleMention1[]', value: 'Role1' },
            { name: 'newRoleMention1[]', value: 'Role2' },
            { name: 'newNgOrWord1[]', value: 'NGWord1' },
            { name: 'newNgOrWord1[]', value: 'NGWord2' },
            { name: 'newNgAndWord1[]', value: 'NGWord3' },
            { name: 'newNgAndWord1[]', value: 'NGWord4' },
            { name: 'newSearchOrWord1[]', value: 'SearchWord1' },
            { name: 'newSearchOrWord1[]', value: 'SearchWord2' },
            { name: 'newSearchAndWord1[]', value: 'SearchWord3' },
            { name: 'newSearchAndWord1[]', value: 'SearchWord4' },
            // 他のフィールドも同様に追加
        ];
        const formData = new FormData();
        formElements.forEach(el => formData.append(el.name, el.value));

        // モックのFormDataを模倣するための関数
        formData.get = (key) => {
            return formElements.find(el => el.name === key)?.value || null;
        };
        formData.getAll = (key) => {
            return formElements.filter(el => el.name === key).map(el => el.value);
        };

        // 関数の実行
        const result = await createJsonData(formElements, formData);

        // 期待値との比較
        expect(result).toBeDefined();
        expect(JSON.parse(result)).toEqual({
            newWebhooks: [
                {
                    webhookId: 'Type1',
                    subscriptionType: 'Subscription1',
                    subscriptionId: 'SubId1',
                    mentionRoles: ['Role1', 'Role2'],
                    mentionUsers: ['Member1', 'Member2'],
                    ngOrWords: ['NGWord1', 'NGWord2'],
                    ngAndWords: ['NGWord3', 'NGWord4'],
                    searchOrWords: ['SearchWord1', 'SearchWord2'],
                    searchAndWords: ['SearchWord3', 'SearchWord4'],
                    mentionOrWords: [],
                    mentionAndWords: [],
                }
            ],
            updateWebhooks: [] // updateWebhooksの期待値も同様に設定
        });
    });

    it('should correctly process form data into JSON with multiple mentions, NG words, search words, and mention words', async () => {
        // モックデータの準備
        const formElements = [
            { name: 'newWebhookType1', value: 'Type1' },
            { name: 'newSubscriptionType1', value: 'Subscription1' },
            { name: 'newSubscriptionId1', value: 'SubId1' },
            { name: 'newMemberMention1[]', value: 'Member1' },
            { name: 'newMemberMention1[]', value: 'Member2' },
            { name: 'newRoleMention1[]', value: 'Role1' },
            { name: 'newRoleMention1[]', value: 'Role2' },
            { name: 'newNgOrWord1[]', value: 'NGWord1' },
            { name: 'newNgOrWord1[]', value: 'NGWord2' },
            { name: 'newNgAndWord1[]', value: 'NGWord3' },
            { name: 'newNgAndWord1[]', value: 'NGWord4' },
            { name: 'newSearchOrWord1[]', value: 'SearchWord1' },
            { name: 'newSearchOrWord1[]', value: 'SearchWord2' },
            { name: 'newSearchAndWord1[]', value: 'SearchWord3' },
            { name: 'newSearchAndWord1[]', value: 'SearchWord4' },
            { name: 'newMentionOrWord1[]', value: 'MentionWord1' },
            { name: 'newMentionOrWord1[]', value: 'MentionWord2' },
            { name: 'newMentionAndWord1[]', value: 'MentionWord3' },
            { name: 'newMentionAndWord1[]', value: 'MentionWord4' },
            // 他のフィールドも同様に追加
        ];
        const formData = new FormData();
        formElements.forEach(el => formData.append(el.name, el.value));

        // モックのFormDataを模倣するための関数
        formData.get = (key) => {
            return formElements.find(el => el.name === key)?.value || null;
        }
        formData.getAll = (key) => {
            return formElements.filter(el => el.name === key).map(el => el.value);
        }

        // 関数の実行
        const result = await createJsonData(formElements, formData);

        // 期待値との比較
        expect(result).toBeDefined();
        expect(JSON.parse(result)).toEqual({
            newWebhooks: [
                {
                    webhookId: 'Type1',
                    subscriptionType: 'Subscription1',
                    subscriptionId: 'SubId1',
                    mentionRoles: ['Role1', 'Role2'],
                    mentionUsers: ['Member1', 'Member2'],
                    ngOrWords: ['NGWord1', 'NGWord2'],
                    ngAndWords: ['NGWord3', 'NGWord4'],
                    searchOrWords: ['SearchWord1', 'SearchWord2'],
                    searchAndWords: ['SearchWord3', 'SearchWord4'],
                    mentionOrWords: ['MentionWord1', 'MentionWord2'],
                    mentionAndWords: ['MentionWord3', 'MentionWord4'],
                }
            ],
            updateWebhooks: [] // updateWebhooksの期待値も同様に設定
        });
    });

    it('should correctly process form data into JSON with multiple mentions, NG words, search words, mention words, and update webhooks', async () => {
        // モックデータの準備
        const formElements = [
            { name: 'newWebhookType1', value: 'Type1' },
            { name: 'newSubscriptionType1', value: 'Subscription1' },
            { name: 'newSubscriptionId1', value: 'SubId1' },
            { name: 'newMemberMention1[]', value: 'Member1' },
            { name: 'newMemberMention1[]', value: 'Member2' },
            { name: 'newRoleMention1[]', value: 'Role1' },
            { name: 'newRoleMention1[]', value: 'Role2' },
            { name: 'newNgOrWord1[]', value: 'NGWord1' },
            { name: 'newNgOrWord1[]', value: 'NGWord2' },
            { name: 'newNgAndWord1[]', value: 'NGWord3' },
            { name: 'newNgAndWord1[]', value: 'NGWord4' },
            { name: 'newSearchOrWord1[]', value: 'SearchWord1' },
            { name: 'newSearchOrWord1[]', value: 'SearchWord2' },
            { name: 'newSearchAndWord1[]', value: 'SearchWord3' },
            { name: 'newSearchAndWord1[]', value: 'SearchWord4' },
            { name: 'newMentionOrWord1[]', value: 'MentionWord1' },
            { name: 'newMentionOrWord1[]', value: 'MentionWord2' },
            { name: 'newMentionAndWord1[]', value: 'MentionWord3' },
            { name: 'newMentionAndWord1[]', value: 'MentionWord4' },
            { name: 'webhookType1', value: 'UpdateType1' },
            { name: 'subscriptionType1', value: 'UpdateSubscription1' },
            { name: 'subscriptionId1', value: 'UpdateSubId1' },
            { name: 'updateMemberMention1[]', value: 'UpdateMember1' },
            { name: 'updateMemberMention1[]', value: 'UpdateMember2' },
            { name: 'updateRoleMention1[]', value: 'UpdateRole1' },
            { name: 'updateRoleMention1[]', value: 'UpdateRole2' },
            { name: 'updateNgOrWord1[]', value: 'UpdateNGWord1' },
            { name: 'updateNgOrWord1[]', value: 'UpdateNGWord2' },
            { name: 'updateNgAndWord1[]', value: 'UpdateNGWord3' },
            { name: 'updateNgAndWord1[]', value: 'UpdateNGWord4' },
            { name: 'updateSearchOrWord1[]', value: 'UpdateSearchWord1' },
            { name: 'updateSearchOrWord1[]', value: 'UpdateSearchWord2' },
            { name: 'updateSearchAndWord1[]', value: 'UpdateSearchWord3' },
            { name: 'updateSearchAndWord1[]', value: 'UpdateSearchWord4' },
            { name: 'updateMentionOrWord1[]', value: 'UpdateMentionWord1' },
            { name: 'updateMentionOrWord1[]', value: 'UpdateMentionWord2' },
            { name: 'updateMentionAndWord1[]', value: 'UpdateMentionWord3' },
            { name: 'updateMentionAndWord1[]', value: 'UpdateMentionWord4' },
            // 他のフィールドも同様に追加
        ];
        const formData = new FormData();
        formElements.forEach(el => formData.append(el.name, el.value));

        // モックのFormDataを模倣するための関数
        formData.get = (key) => {
            return formElements.find(el => el.name === key)?.value || null;
        }
        formData.getAll = (key) => {
            return formElements.filter(el => el.name === key).map(el => el.value);
        }

        // 関数の実行
        const result = await createJsonData(formElements, formData);

        // 期待値との比較
        expect(result).toBeDefined();
        expect(JSON.parse(result)).toEqual({
            newWebhooks: [
                {
                    webhookId: 'Type1',
                    subscriptionType: 'Subscription1',
                    subscriptionId: 'SubId1',
                    mentionRoles: ['Role1', 'Role2'],
                    mentionUsers: ['Member1', 'Member2'],
                    ngOrWords: ['NGWord1', 'NGWord2'],
                    ngAndWords: ['NGWord3', 'NGWord4'],
                    searchOrWords: ['SearchWord1', 'SearchWord2'],
                    searchAndWords: ['SearchWord3', 'SearchWord4'],
                    mentionOrWords: ['MentionWord1', 'MentionWord2'],
                    mentionAndWords: ['MentionWord3', 'MentionWord4'],
                }
            ],
            updateWebhooks: [
                {
                    webhookSerialId: 1,
                    webhookId: 'UpdateType1',
                    subscriptionType: 'UpdateSubscription1',
                    subscriptionId: 'UpdateSubId1',
                    mentionRoles: ['UpdateRole1', 'UpdateRole2'],
                    mentionUsers: ['UpdateMember1', 'UpdateMember2'],
                    ngOrWords: ['UpdateNGWord1', 'UpdateNGWord2'],
                    ngAndWords: ['UpdateNGWord3', 'UpdateNGWord4'],
                    searchOrWords: ['UpdateSearchWord1', 'UpdateSearchWord2'],
                    searchAndWords: ['UpdateSearchWord3', 'UpdateSearchWord4'],
                    mentionOrWords: ['UpdateMentionWord1', 'UpdateMentionWord2'],
                    mentionAndWords: ['UpdateMentionWord3', 'UpdateMentionWord4'],
                    deleteFlag: false,
                }
            ]
        });
    });
});
