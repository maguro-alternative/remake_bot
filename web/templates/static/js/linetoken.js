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

    console.log(jsonData);

    // データを送信
    await fetch(`/api/${guildId}/linetoken`, {
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
    data['debugMode'] === 'on' ? data['debugMode'] = true : data['debugMode'] = false;
    data['lineNotifyTokenDelete'] === 'on' ? data['lineNotifyTokenDelete'] = true : data['lineNotifyTokenDelete'] = false;
    data['lineBotTokenDelete'] === 'on' ? data['lineBotTokenDelete'] = true : data['lineBotTokenDelete'] = false;
    data['lineBotSecretDelete'] === 'on' ? data['lineBotSecretDelete'] = true : data['lineBotSecretDelete'] = false;
    data['lineGroupIdDelete'] === 'on' ? data['lineGroupIdDelete'] = true : data['lineGroupIdDelete'] = false;
    data['lineClientIdDelete'] === 'on' ? data['lineClientIdDelete'] = true : data['lineClientIdDelete'] = false;
    data['lineClientSecretDelete'] === 'on' ? data['lineClientSecretDelete'] = true : data['lineClientSecretDelete'] = false;
    return JSON.stringify(data);
}

try {
    module.exports = { createJsonData };
} catch (e) {
}
