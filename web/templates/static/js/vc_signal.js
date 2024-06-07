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
    await fetch(`/api/${guildId}/vc-signal`, {
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

const createJsonData = async function(formElements, formData) {
    // それぞれ入力内容を配列に格納
    let ngUserArray, ngRoleArray, mentionUserArray, mentionRoleArray;
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
                "vcChannelId": channelId,
                "sendSignal": false,
                "sendChannelId": formData.get(`defaultChannelId${channelId}`),
                "joinBot": false,
                "everyoneMention": false,
                "vcSignalMentionUsers": mentionUserArray,
                "vcSignalMentionRoles": mentionRoleArray,
                "vcSignalNgUsers": ngUserArray,
                "vcSignalNgRoles": ngRoleArray
            }
        }
        if ((formKey.includes('vcSignalMentionUserIds')) && (formKey.includes('[]'))) {
            jsonTmp[channelId]['vcSignalMentionUserIds'] = formData.getAll(formKey);
        }
        if ((formKey.includes('vcSignalMentionRoleIds')) && (formKey.includes('[]'))) {
            jsonTmp[channelId]['vcSignalMentionRoleIds'] = formData.getAll(formKey);
        }
        if ((formKey.includes('vcSignalNgUserIds')) && (formKey.includes('[]'))) {
            jsonTmp[channelId]['vcSignalNgUserIds'] = formData.getAll(formKey);
        }
        if ((formKey.includes('vcSignalNgRoleIds')) && (formKey.includes('[]'))) {
            jsonTmp[channelId]['vcSignalNgRoleIds'] = formData.getAll(formKey);
        }
        if (formKey.includes('sendSignal')) {
            document.getElementById(formKey).checked ? jsonTmp[channelId]["sendSignal"] = true : jsonTmp[channelId]["sendSignal"] = false;
        }
        if (formKey.includes('defaultChannelId')) {
            jsonTmp[channelId]["defaultChannelId"] = formData.get(formKey);
        }
        if (formKey.includes('joinBot')) {
            document.getElementById(formKey).checked ? jsonTmp[channelId]["joinBot"] = true : jsonTmp[channelId]["joinBot"] = false;
        }
        if (formKey.includes('everyoneMention')) {
            document.getElementById(formKey).checked ? jsonTmp[channelId]["everyoneMention"] = true : jsonTmp[channelId]["everyoneMention"] = false;
        }
    }
    for (let key in jsonTmp) {
        channelArray.push(jsonTmp[key]);
    }

    return JSON.stringify({
        "vcSignals": channelArray
    });
}

try {
    module.exports = { createJsonData };
} catch (e) {
}
