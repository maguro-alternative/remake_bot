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
    const data = await fetchGroupData(guildId, jsonData);
    if (data.ok) {
        alert('設定を保存しました');
        window.location.href = `/guild/${guildId}`;
    } else {
        window.alert('送信に失敗しました');
    }

}

const createJsonData = async function(formData) {
    const data = Object.fromEntries(formData.entries());

    return JSON.stringify(data);
}

const fetchGroupData = async function(guildId, jsonData) {
    const res = await fetch(`/api/${guildId}/group`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: jsonData
    })

    return res;
}

try {
    module.exports = { fetchGroupData, createJsonData };
} catch (e) {
}
