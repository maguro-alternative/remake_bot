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
    // それぞれ入力内容を配列に格納

    const jsonData = await createJsonData(formElements, formData);

    // データを送信
    await fetch(`/api/${guildId}/line-post-discord-channel`, {
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
    })
}

async function createJsonData(formElements, formData) {
    // それぞれ入力内容を配列に格納
    let ngTypeArray, ngUserArray, ngRoleArray;
    const channelArray = [];
    let jsonTmp = {};
    for (let i = 0; i < formElements.length; i++) {
        formKey = formElements[i].name;
        channelIdmatch = formKey.match(/[0-9]/g);
        if (channelIdmatch === null) {
            continue;
        }
        channelId = channelIdmatch.join('');
        if (jsonTmp[channelId] === undefined) {
            jsonTmp[channelId] = {
                "channel_id": channelId,
                "ng": false,
                "bot_message": false,
                "ng_types": ngTypeArray,
                "ng_users": ngUserArray,
                "ng_roles": ngRoleArray
            }
        }
        if ((formKey.includes('ng_types')) && (formKey.includes('[]'))) {
            jsonTmp[channelId]['ng_types'] = formData.getAll(formKey).map( str => parseInt(str, 10) );
        }
        if ((formKey.includes('ng_users')) && (formKey.includes('[]'))) {
            jsonTmp[channelId]['ng_users'] = formData.getAll(formKey);
        }
        if ((formKey.includes('ng_roles')) && (formKey.includes('[]'))) {
            jsonTmp[channelId]['ng_roles'] = formData.getAll(formKey);
        }
        if ((formKey.includes('ng_')) && (!formKey.includes('[]'))) {
            document.getElementById(formKey).checked ? jsonTmp[channelId]["ng"] = true : jsonTmp[channelId]["ng"] = false;
        }
        if (formKey.includes('bot_message')) {
            document.getElementById(formKey).checked ? jsonTmp[channelId]["bot_message"] = true : jsonTmp[channelId]["bot_message"] = false;
        }
    }

    for (let key in jsonTmp) {
        channelArray.push(jsonTmp[key]);
    }

    return JSON.stringify({
        "channels": channelArray
    });
}
