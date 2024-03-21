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
    let permissionName;
    // 各formのkeyを取得
    for (let i = 0; i < formElements.length; i++) {
        formKey = formElements[i].name;
        if (formKey.includes('permission_code')) {
            permissionName = formKey.substr(0, str.indexOf('_permission_code'));
        } else if (formKey.includes('member_permission_id')) {
            permissionName = formKey.substr(0, str.indexOf('_member_permission_id'));
            formData.getAll(formKey);
        } else if (formKey.includes('role_permission_id')) {
            permissionName = formKey.substr(0, str.indexOf('_role_permission_id'));
            formData.getAll(formKey);
        }
    }
    const data = Object.fromEntries(formData.entries());

    const jsonData = JSON.stringify(data);

    // データを送信
    await fetch(`/api/${guildId}/permission`, {
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