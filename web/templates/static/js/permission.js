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
    let jsonTmp = {
        "permission_codes": [],
        "permission_ids": [],
    };
    console.log(formElements)
    // 各formのkeyを取得
    for (let i = 0; i < formElements.length; i++) {
        formKey = formElements[i].name;
        console.log(formKey);
        if (formKey.includes('permission_code')) {
            permissionName = formKey.substr(0, formKey.indexOf('_permission_code'));
            jsonTmp['permission_codes'].push({
                "guild_id": guildId,
                "type": permissionName,
                "code": parseInt(formData.get(formKey))
            })
        } else if (formKey.includes('member_permission_id')) {
            permissionName = formKey.substr(0, formKey.indexOf('_member_permission_id'));
            for (let user of formData.getAll(formKey)) {
                console.log(user);
                jsonTmp['permission_ids'].push({
                    "guild_id": guildId,
                    "type": permissionName,
                    "target_type": "user",
                    "target_id": user,
                    "permission": "all"
                })
            }
        } else if (formKey.includes('role_permission_id')) {
            permissionName = formKey.substr(0, formKey.indexOf('_role_permission_id'));
            for (let role of formData.getAll(formKey)) {
                jsonTmp['permission_ids'].push({
                    "guild_id": guildId,
                    "type": permissionName,
                    "target_type": "role",
                    "target_id": role,
                    "permission": "all"
                })
            }
        }
    }

    const jsonData = JSON.stringify(jsonTmp);

    console.log(jsonData)

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