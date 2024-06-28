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

    // 他のテストケースも同様に追加
});
