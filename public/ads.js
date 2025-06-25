window.onload = function() {
    setTimeout(function () {
        document.querySelector('.popup').classList.add('visible');
        document.querySelector('.overlay').classList.add('visible');
        let progressBar = document.getElementById('progressBar');
        let width = 0;
        let interval = setInterval(function () {
            if (width >= 100) {
                clearInterval(interval);
                document.getElementById('closeBtn').style.display = 'block';
            } else {
                width++;
                progressBar.style.width = width + '%';
            }
        }, 50);
    }, 5000);
    document.getElementById('closeBtn').addEventListener('click', function () {
        document.querySelector('.popup').classList.remove('visible');
        document.querySelector('.overlay').classList.remove('visible');
    });
};
