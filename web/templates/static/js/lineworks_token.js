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
  await fetch(`/api/${guildId}/lineworks-token`, {
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
  data['lineWorksClientIdDelete'] === 'on' ? data['lineWorksClientIdDelete'] = true : data['lineWorksClientIdDelete'] = false;
  data['lineWorksClientSecretDelete'] === 'on' ? data['lineWorksClientSecretDelete'] = true : data['lineWorksClientSecretDelete'] = false;
  data['lineWorksServiceAccountDelete'] === 'on' ? data['lineWorksServiceAccountDelete'] = true : data['lineWorksServiceAccountDelete'] = false;
  data['lineWorksPrivateKeyDelete'] === 'on' ? data['lineWorksPrivateKeyDelete'] = true : data['lineWorksPrivateKeyDelete'] = false;
  data['lineWorksDomainIdDelete'] === 'on' ? data['lineWorksDomainIdDelete'] = true : data['lineWorksDomainIdDelete'] = false;
  data['lineWorksAdminIdDelete'] === 'on' ? data['lineWorksAdminIdDelete'] = true : data['lineWorksAdminIdDelete'] = false;
  data['lineWorksBotIdDelete'] === 'on' ? data['lineWorksBotIdDelete'] = true : data['lineWorksBotIdDelete'] = false;
  data['lineWorksBotSecretDelete'] === 'on' ? data['lineWorksBotSecretDelete'] = true : data['lineWorksBotSecretDelete'] = false;
  data['lineWorksGroupIdDelete'] === 'on' ? data['lineWorksGroupIdDelete'] = true : data['lineWorksGroupIdDelete'] = false;
  return JSON.stringify(data);
}

try {
  module.exports = { createJsonData };
} catch (e) {
}
