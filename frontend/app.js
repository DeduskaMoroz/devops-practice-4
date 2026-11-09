fetch('/api/items')
  .then(res => res.json())
  .then(data => {
    const list = document.getElementById('items-list');
    list.innerHTML = data.map(item => `<li>#${item.id}: ${item.name}</li>`).join('');
  })
  .catch(err => {
    document.getElementById('items-list').innerText = 'Ошибка загрузки: ' + err;
  });