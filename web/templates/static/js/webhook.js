// submit時にイベント実行をする関数
document.getElementById('form').onsubmit = async function (event) {
    let check = window.confirm('送信します。よろしいですか？');
    if (check === false) {
        event.preventDefault();
        return;
    }
    // 再読み込み防止
    event.preventDefault();
    const guildId = location.pathname.match(/[0-9]/g).join('');
    const formData = new FormData(document.getElementById('form'));
    const formElements = document.forms['form'].elements;

    const jsonData = await createJsonData(formElements, formData);

    // データを送信
    await fetch(`/api/${guildId}/webhook`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: jsonData
    }).then((res) => {
        if (res.ok && res.status === 200) {
            alert('設定を保存しました');
            window.location.href = `/guild/${guildId}`;
        } else {
            alert('設定の保存に失敗しました');
        }
    });
}

// JavaScriptでsessionStorageにoptionを保存
window.onload = function() {
    const webhookOptions = document.querySelector('#newWebhookType1').innerHTML;
    window.sessionStorage.setItem('webhookOptions', webhookOptions);

    const subscriptionNameOptions = document.querySelector('#newSubscriptionName1').innerHTML;
    window.sessionStorage.setItem("subscriptionNameOptions", subscriptionNameOptions);

    // IDに[]が含まれる場合、\\[と\\]でエスケープする
    const memberMentionOptions = document.querySelector('#newMemberMention1\\[\\]').innerHTML;
    window.sessionStorage.setItem('memberMentionOptions', memberMentionOptions);

    const roleMentionOptions = document.querySelector('#newRoleMention1\\[\\]').innerHTML;
    window.sessionStorage.setItem('roleMentionOptions', roleMentionOptions);
}

const createJsonData = async function(formElements, formData) {
    let newWebhooksTmp = {}, updateWebhooksTmp = {};
    for (let i = 0; i < formElements.length; i++) {
        formKey = formElements[i].name;
        if (formKey === '') {
            continue;
        }
        webhookFormId = formKey.match(/[0-9]/g).join('');
        if (newWebhooksTmp[webhookFormId] === undefined && formKey.includes('newWebhookType')) {
            newWebhooksTmp[webhookFormId] = {
                "webhookId":"",
                "subscriptionType":"",
                "subscriptionId":"",
                "mentionRoles":[],
                "mentionUsers":[],
                "ngOrWords":[],
                "ngAndWords":[],
                "searchOrWords":[],
                "searchAndWords":[],
                "mentionOrWords":[],
                "mentionAndWords":[],
            }
        }
        if (updateWebhooksTmp[webhookFormId] === undefined && formKey.includes('updateWebhookType')) {
            updateWebhooksTmp[webhookFormId] = {
                "webhookSerialId":Number(webhookFormId),
                "webhookId":"",
                "subscriptionType":"",
                "subscriptionId":"",
                "mentionRoles":[],
                "mentionUsers":[],
                "ngOrWords":[],
                "ngAndWords":[],
                "searchOrWords":[],
                "searchAndWords":[],
                "mentionOrWords":[],
                "mentionAndWords":[],
                "deleteFlag":false
            }
        }
        if (formKey.includes('newWebhookType')) {
            newWebhooksTmp[webhookFormId]['webhookId'] = formData.get(formKey);
        }
        if (formKey.includes('newSubscriptionName')) {
            newWebhooksTmp[webhookFormId]['subscriptionType'] = formData.get(formKey);
        }
        if (formKey.includes('newSubscriptionId')) {
            newWebhooksTmp[webhookFormId]['subscriptionId'] = formData.get(formKey);
        }
        if (formKey.includes('newMemberMention')) {
            newWebhooksTmp[webhookFormId]['mentionUsers'] = formData.getAll(formKey);
        }
        if (formKey.includes('newRoleMention')) {
            newWebhooksTmp[webhookFormId]['mentionRoles'] = formData.getAll(formKey);
        }
        if (formKey.includes('newNgOr')) {
            newWebhooksTmp[webhookFormId]['ngOrWords'] = formData.getAll(formKey);
        }
        if (formKey.includes('newNgAnd')) {
            newWebhooksTmp[webhookFormId]['ngAndWords'] = formData.getAll(formKey);
        }
        if (formKey.includes('newSearchOr')) {
            newWebhooksTmp[webhookFormId]['searchOrWords'] = formData.getAll(formKey);
        }
        if (formKey.includes('newSearchAnd')) {
            newWebhooksTmp[webhookFormId]['searchAndWords'] = formData.getAll(formKey);
        }
        if (formKey.includes('newMentionOr')) {
            newWebhooksTmp[webhookFormId]['mentionOrWords'] = formData.getAll(formKey);
        }
        if (formKey.includes('newMentionAnd')) {
            newWebhooksTmp[webhookFormId]['mentionAndWords'] = formData.getAll(formKey);
        }
        if (formKey.includes('updateWebhookType')) {
            updateWebhooksTmp[webhookFormId]['webhookId'] = formData.get(formKey);
        }
        if (formKey.includes('updateSubscriptionName')) {
            updateWebhooksTmp[webhookFormId]['subscriptionType'] = formData.get(formKey);
        }
        if (formKey.includes('updateSubscriptionId')) {
            updateWebhooksTmp[webhookFormId]['subscriptionId'] = formData.get(formKey);
        }
        if (formKey.includes('updateMemberMention')) {
            updateWebhooksTmp[webhookFormId]['mentionUsers'] = formData.getAll(formKey);
        }
        if (formKey.includes('updateRoleMention')) {
            updateWebhooksTmp[webhookFormId]['mentionRoles'] = formData.getAll(formKey);
        }
        if (formKey.includes('updateNgOr')) {
            updateWebhooksTmp[webhookFormId]['ngOrWords'] = formData.getAll(formKey);
        }
        if (formKey.includes('updateNgAnd')) {
            updateWebhooksTmp[webhookFormId]['ngAndWords'] = formData.getAll(formKey);
        }
        if (formKey.includes('updateSearchOr')) {
            updateWebhooksTmp[webhookFormId]['searchOrWords'] = formData.getAll(formKey);
        }
        if (formKey.includes('updateSearchAnd')) {
            updateWebhooksTmp[webhookFormId]['searchAndWords'] = formData.getAll(formKey);
        }
        if (formKey.includes('updateMentionOr')) {
            updateWebhooksTmp[webhookFormId]['mentionOrWords'] = formData.getAll(formKey);
        }
        if (formKey.includes('updateMentionAnd')) {
            updateWebhooksTmp[webhookFormId]['mentionAndWords'] = formData.getAll(formKey);
        }
        if (formKey.includes('updateDeleteFlag')) {
            document.getElementById(formKey).checked ? updateWebhooksTmp[webhookFormId]["deleteFlag"] = true : updateWebhooksTmp[webhookFormId]["deleteFlag"] = false;
        }
    }
    const newWebhooks = Object.values(newWebhooksTmp);
    const updateWebhooks = Object.values(updateWebhooksTmp);
    console.log(JSON.stringify({
        "newWebhooks": newWebhooks,
        "updateWebhooks": updateWebhooks
    }));
    return JSON.stringify({
        "newWebhooks": newWebhooks,
        "updateWebhooks": updateWebhooks
    });
}

/*function addWebhook() {
    const webhook = document.getElementById('new_webhook');
    const webhookClone = webhook.cloneNode(true);
    webhookClone.id = `new_webhook${webhook.children.length}`;
    webhookClone.querySelector('input').value = '';
    webhook.appendChild(webhookClone);
}*/

// 追加ボタンがクリックされたときに実行される関数
function addWebhook() {
    const webhook = document.getElementById('newWebhook');
    // 現在のWebhook要素の数を取得
    const currentWebhooks = document.querySelectorAll('[id^="newWebhookType"]').length;
    // 新しいIDの番号を設定
    const newIdNumber = currentWebhooks + 1;

    while (document.getElementById(`newWebhookType${newIdNumber}`) !== null) {
        newIdNumber++;
    }

    // 新しいWebhook要素を作成
    const newWebhookLabel = document.createElement('label');
    newWebhookLabel.htmlFor = `newWebhookType${newIdNumber}`;
    newWebhookLabel.textContent = 'Webhook';
    const newWebhookType = document.createElement('select');
    newWebhookType.name = `newWebhookType${newIdNumber}`;
    newWebhookType.id = `newWebhookType${newIdNumber}`;
    newWebhookType.innerHTML = sessionStorage.getItem('webhookOptions');

    const newSubscriptionNameLabel = document.createElement('label');
    newSubscriptionNameLabel.htmlFor = `newSubscriptionName${newIdNumber}`;
    newSubscriptionNameLabel.textContent = 'サービス名';
    const newSubscriptionName = document.createElement('select');
    newSubscriptionName.name = `newSubscriptionName${newIdNumber}`;
    newSubscriptionName.innerHTML = sessionStorage.getItem("subscriptionNameOptions")

    const newSubscriptionIdLabel = document.createElement('label');
    newSubscriptionIdLabel.htmlFor = `newSubscriptionId${newIdNumber}`;
    newSubscriptionIdLabel.textContent = 'サービスID';
    const newSubscriptionId = document.createElement('input');
    newSubscriptionId.type = 'text';
    newSubscriptionId.name = `newSubscriptionId${newIdNumber}`;

    const newMenberMentionLabel = document.createElement('label');
    newMenberMentionLabel.htmlFor = `newMemberMention${newIdNumber}[]`;
    newMenberMentionLabel.textContent = 'メンションするユーザー';
    const newMenberMention = document.createElement('select');
    newMenberMention.name = `newMemberMention${newIdNumber}[]`;
    newMenberMention.innerHTML = sessionStorage.getItem('memberMentionOptions');
    newMenberMention.multiple = true;

    const newRoleMentionLabel = document.createElement('label');
    newRoleMentionLabel.htmlFor = `newRoleMention${newIdNumber}[]`;
    newRoleMentionLabel.textContent = 'メンションするロール';
    const newRoleMention = document.createElement('select');
    newRoleMention.name = `newRoleMention${newIdNumber}[]`;
    newRoleMention.innerHTML = sessionStorage.getItem('roleMentionOptions');
    newRoleMention.multiple = true;

    const newNgOrWordsDiv = document.createElement('div');
    newNgOrWordsDiv.id = `newNgOrWords${newIdNumber}`;
    const newNgOrWords = document.createElement('button');
    newNgOrWords.type = 'button';
    newNgOrWords.textContent = 'NGワードOR追加';
    newNgOrWords.onclick = function () {
        addWord('newNgOr', newIdNumber);
    };
    newNgOrWordsDiv.appendChild(newNgOrWords);

    const newNgAndWordsDiv = document.createElement('div');
    newNgAndWordsDiv.id = `newNgAndWords${newIdNumber}`;
    const newNgAndWords = document.createElement('button');
    newNgAndWords.type = 'button';
    newNgAndWords.textContent = 'NGワードAND追加';
    newNgAndWords.onclick = function () {
        addWord('newNgAnd', newIdNumber);
    };
    newNgAndWordsDiv.appendChild(newNgAndWords);

    const newSearchOrWordsDiv = document.createElement('div');
    newSearchOrWordsDiv.id = `newSearchOrWords${newIdNumber}`;
    const newSearchOrWords = document.createElement('button');
    newSearchOrWords.type = 'button';
    newSearchOrWords.textContent = '検索ワードOR追加';
    newSearchOrWords.onclick = function () {
        addWord('newSearchOr', newIdNumber);
    };
    newSearchOrWordsDiv.appendChild(newSearchOrWords);

    const newSearchAndWordsDiv = document.createElement('div');
    newSearchAndWordsDiv.id = `newSearchAndWords${newIdNumber}`;
    const newSearchAndWords = document.createElement('button');
    newSearchAndWords.type = 'button';
    newSearchAndWords.textContent = '検索ワードAND追加';
    newSearchAndWords.onclick = function () {
        addWord('newSearchAnd', newIdNumber);
    };
    newSearchAndWordsDiv.appendChild(newSearchAndWords);

    const newMentionOrWordsDiv = document.createElement('div');
    newMentionOrWordsDiv.id = `newMentionOrWords${newIdNumber}`;
    const newMentionOrWords = document.createElement('button');
    newMentionOrWords.type = 'button';
    newMentionOrWords.textContent = 'メンションOR追加';
    newMentionOrWords.onclick = function () {
        addWord('newMentionOr', newIdNumber);
    };
    newMentionOrWordsDiv.appendChild(newMentionOrWords);

    const newMentionAndWordsDiv = document.createElement('div');
    newMentionAndWordsDiv.id = `newMentionAndWords${newIdNumber}`;
    const newMentionAndWords = document.createElement('button');
    newMentionAndWords.type = 'button';
    newMentionAndWords.textContent = 'メンションAND追加';
    newMentionAndWords.onclick = function () {
        addWord('newMentionAnd', newIdNumber);
    };
    newMentionAndWordsDiv.appendChild(newMentionAndWords);

    // 新しいWebhook要素を追加
    webhook.appendChild(document.createElement('br'));
    webhook.appendChild(document.createElement('br'));
    webhook.appendChild(newWebhookLabel);
    webhook.appendChild(newWebhookType);
    webhook.appendChild(document.createElement('br'));
    webhook.appendChild(document.createElement('br'));
    webhook.appendChild(newSubscriptionNameLabel);
    webhook.appendChild(newSubscriptionName);
    webhook.appendChild(document.createElement('br'));
    webhook.appendChild(document.createElement('br'));
    webhook.appendChild(newSubscriptionIdLabel);
    webhook.appendChild(newSubscriptionId);
    webhook.appendChild(document.createElement('br'));
    webhook.appendChild(document.createElement('br'));
    webhook.appendChild(newMenberMentionLabel);
    webhook.appendChild(newMenberMention);
    webhook.appendChild(document.createElement('br'));
    webhook.appendChild(document.createElement('br'));
    webhook.appendChild(newRoleMentionLabel);
    webhook.appendChild(newRoleMention);
    webhook.appendChild(document.createElement('br'));
    webhook.appendChild(document.createElement('br'));
    webhook.appendChild(newNgOrWordsDiv);
    webhook.appendChild(document.createElement('br'));
    webhook.appendChild(document.createElement('br'));
    webhook.appendChild(newNgAndWordsDiv);
    webhook.appendChild(document.createElement('br'));
    webhook.appendChild(document.createElement('br'));
    webhook.appendChild(newSearchOrWordsDiv);
    webhook.appendChild(document.createElement('br'));
    webhook.appendChild(document.createElement('br'));
    webhook.appendChild(newSearchAndWordsDiv);
    webhook.appendChild(document.createElement('br'));
    webhook.appendChild(document.createElement('br'));
    webhook.appendChild(newMentionOrWordsDiv);
    webhook.appendChild(document.createElement('br'));
    webhook.appendChild(document.createElement('br'));
    webhook.appendChild(newMentionAndWordsDiv);
}

function addWord(condition, serialId) {
    const word = document.getElementById(`${condition}Words${serialId}`);
    const wordClone = word.cloneNode(true);
    wordClone.id = `word${serialId}`;

    // inputタグが存在しないため、新しいinput要素を作成して追加します。
    const input = document.createElement('input');
    input.type = 'text';
    input.value = ''; // 初期値は空文字列
    input.name = `${condition}Word${serialId}[]`;

    // 削除ボタンを作成
    const deleteButton = document.createElement('button');
    deleteButton.textContent = '削除';

    // 削除ボタンのクリックイベントリスナーを追加
    deleteButton.addEventListener('click', function() {
        input.remove(); // input要素を削除
        deleteButton.remove(); // 削除ボタン自身も削除
    });
    // input要素と削除ボタンを親要素に追加
    word.appendChild(input);
    word.appendChild(deleteButton);
}

try {
    module.exports = { createJsonData, addWebhook, addWord};
} catch (e) {
}
