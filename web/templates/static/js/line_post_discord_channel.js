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

const createJsonData = async function(formElements, formData) {
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
                "channelId": channelId,
                "ng": false,
                "botMessage": false,
                "ngTypes": ngTypeArray,
                "ngUsers": ngUserArray,
                "ngRoles": ngRoleArray
            }
        }
        if ((formKey.includes('ngTypes')) && (formKey.includes('[]'))) {
            jsonTmp[channelId]['ngTypes'] = formData.getAll(formKey).map( str => parseInt(str, 10) );
        }
        if ((formKey.includes('ngUsers')) && (formKey.includes('[]'))) {
            jsonTmp[channelId]['ngUsers'] = formData.getAll(formKey);
        }
        if ((formKey.includes('ngRoles')) && (formKey.includes('[]'))) {
            jsonTmp[channelId]['ngRoles'] = formData.getAll(formKey);
        }
        if ((formKey.includes('ng')) && (!formKey.includes('[]'))) {
            document.getElementById(formKey).checked ? jsonTmp[channelId]["ng"] = true : jsonTmp[channelId]["ng"] = false;
        }
        if (formKey.includes('botMessage')) {
            document.getElementById(formKey).checked ? jsonTmp[channelId]["botMessage"] = true : jsonTmp[channelId]["botMessage"] = false;
        }
    }

    for (let key in jsonTmp) {
        channelArray.push(jsonTmp[key]);
    }

    return JSON.stringify({
        "channels": channelArray
    });
}

try {
    module.exports = { createJsonData };
} catch (e) {
}
