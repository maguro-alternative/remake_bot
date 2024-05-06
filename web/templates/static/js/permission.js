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

    const jsonData = await createJsonData(guildId, formData, formElements);

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

const createJsonData = async function(guildId, formData, formElements) {
    let permissionName;
    let jsonTmp = {
        "permissionCodes": [],
        "permissionUserIds": [],
        "permissionRoleIds": [],
    };

    // 各formのkeyを取得
    for (let i = 0; i < formElements.length; i++) {
        formKey = formElements[i].name;
        if (formKey.includes('PermissionCode')) {
            permissionName = formKey.substr(0, formKey.indexOf('PermissionCode'));
            jsonTmp['permissionCodes'].push({
                "guildId": guildId,
                "type": permissionName,
                "code": parseInt(formData.get(formKey))
            })
        } else if (formKey.includes('MemberPermissionId')) {
            permissionName = formKey.substr(0, formKey.indexOf('MemberPermissionId'));
            for (let user of formData.getAll(formKey)) {
                jsonTmp['permissionUserIds'].push({
                    "guildId": guildId,
                    "type": permissionName,
                    "userId": user,
                    "permission": "all"
                })
            }
        } else if (formKey.includes('RolePermissionId')) {
            permissionName = formKey.substr(0, formKey.indexOf('RolePermissionId'));
            for (let role of formData.getAll(formKey)) {
                jsonTmp['permissionRoleIds'].push({
                    "guildId": guildId,
                    "type": permissionName,
                    "roleId": role,
                    "permission": "all"
                })
            }
        }
    }

    return JSON.stringify(jsonTmp);
}

try {
    module.exports = { createJsonData };
} catch (e) {
}
