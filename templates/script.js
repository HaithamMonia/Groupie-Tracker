document.addEventListener('DOMContentLoaded', function() {
    fetch('/dates')
        .then(response => response.json())
        .then(dates => {
            const datesContainer = document.getElementById('dates-container');
            dates.forEach(date => {
                const dateCard = document.createElement('div');
                dateCard.classList.add('date-card');

                const id = document.createElement('p');
                id.textContent = `ID: ${date.id}`;

                const dataKeys = Object.keys(date.data);
                const dataItems = dataKeys.map(key => {
                    const item = document.createElement('p');
                    item.textContent = `${key}: ${date.data[key]}`;
                    return item;
                });

                dateCard.appendChild(id);
                dataItems.forEach(item => dateCard.appendChild(item));
                datesContainer.appendChild(dateCard);
            });
        })
        .catch(error => console.error('Error fetching dates:', error));
});
