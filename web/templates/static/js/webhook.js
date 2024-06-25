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

    const jsonData = await createJsonData(formData);

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

const createJsonData = async function(formData) {
    const data = Object.fromEntries(formData.entries());
    return JSON.stringify(data);
}

function addWebhook() {
    const webhook = document.getElementById('new_webhook');
    const webhookClone = webhook.cloneNode(true);
    webhookClone.id = `new_webhook${webhook.children.length}`;
    webhookClone.querySelector('input').value = '';
    webhook.appendChild(webhookClone);
}

function addWord(condition, serialId) {
    const word = document.getElementById(`${condition}_words`);
    const wordClone = word.cloneNode(true);
    wordClone.id = `word${serialId}`;

    // inputタグが存在しないため、新しいinput要素を作成して追加します。
    const input = document.createElement('input');
    input.type = 'text';
    input.value = ''; // 初期値は空文字列
    //wordClone.appendChild(input); // 作成したinput要素をwordCloneに追加

    //wordClone.querySelector('select').value = condition;
    word.appendChild(input);
}

try {
    module.exports = { createJsonData, addWebhook, addWord};
} catch (e) {
}
