document.addEventListener('DOMContentLoaded', function () {

    // Переключатель языков
    window.setLanguage = function (lang) {
        document.querySelectorAll('[data-lang]').forEach(el => {
            const text = el.getAttribute(`data-lang-${lang}`);
            if (text) {
                el.textContent = text;
            }
        });
        document.documentElement.lang = lang;
        console.log(`Язык переключен на: ${lang}`);
    };

    // Плавная прокрутка по якорям
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            const targetId = this.getAttribute('href');
            if (targetId === '#') return;

            const targetElement = document.querySelector(targetId);
            if (targetElement) {
                window.scrollTo({
                    top: targetElement.offsetTop - 70,
                    behavior: 'smooth'
                });

                const navbarCollapse = document.querySelector('.navbar-collapse');
                if (navbarCollapse && navbarCollapse.classList.contains('show')) {
                    navbarCollapse.classList.remove('show');
                }
            }
        });
    });

    // Обработчик формы обратной связи
    const contactForm = document.getElementById('contactForm');
    const formStatus = document.getElementById('formStatus');

    if (contactForm) {
        contactForm.addEventListener('submit', function (e) {
            e.preventDefault();
            
            const name = document.getElementById('name').value;
            const email = document.getElementById('email').value;
            const message = document.getElementById('message').value;
            
            // Показываем статус отправки
            if (formStatus) {
                formStatus.style.display = 'block';
                formStatus.className = 'alert alert-info mt-3';
                
                const lang = document.documentElement.lang || 'sr';
                if (lang === 'ru') {
                    formStatus.textContent = 'Отправка сообщения...';
                } else {
                    formStatus.textContent = 'Slanje poruke...';
                }
            }
            
            // Отправляем запрос на API основного сайта
            fetch('https://svetu.rs/api/v1/public/send-email', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    name: name,
                    email: email,
                    message: message,
                    source: 'klimagrad'
                })
            })
            .then(response => response.json())
            .then(data => {
                if (formStatus) {
                    formStatus.style.display = 'block';
                    
                    if (data.success) {
                        formStatus.className = 'alert alert-success mt-3';
                        const lang = document.documentElement.lang || 'sr';
                        if (lang === 'ru') {
                            formStatus.textContent = 'Спасибо за сообщение! Мы свяжемся с вами в ближайшее время.';
                        } else {
                            formStatus.textContent = 'Hvala na poruci! Kontaktiraćemo vas uskoro.';
                        }
                        
                        // Очищаем форму
                        contactForm.reset();
                    } else {
                        formStatus.className = 'alert alert-danger mt-3';
                        const lang = document.documentElement.lang || 'sr';
                        if (lang === 'ru') {
                            formStatus.textContent = 'Произошла ошибка при отправке сообщения. Пожалуйста, попробуйте позже.';
                        } else {
                            formStatus.textContent = 'Došlo je do greške prilikom slanja poruke. Molimo pokušajte ponovo kasnije.';
                        }
                    }
                }
            })
            .catch(error => {
                if (formStatus) {
                    formStatus.style.display = 'block';
                    formStatus.className = 'alert alert-danger mt-3';
                    
                    const lang = document.documentElement.lang || 'sr';
                    if (lang === 'ru') {
                        formStatus.textContent = 'Произошла ошибка при отправке сообщения. Пожалуйста, попробуйте позже.';
                    } else {
                        formStatus.textContent = 'Došlo je do greške prilikom slanja poruke. Molimo pokušajte ponovo kasnije.';
                    }
                }
                console.error('Error:', error);
            });
        });
    }

    // Автоматическая пауза видео при скролле за пределы видимости
    const videos = document.querySelectorAll('video');

    if (videos.length > 0) {
        function isElementInViewport(el) {
            const rect = el.getBoundingClientRect();
            return (
                rect.top >= 0 &&
                rect.left >= 0 &&
                rect.bottom <= (window.innerHeight || document.documentElement.clientHeight) &&
                rect.right <= (window.innerWidth || document.documentElement.clientWidth)
            );
        }

        function handleVideoVisibility() {
            videos.forEach(video => {
                if (isElementInViewport(video)) {
                    if (video.paused && video.dataset.autoplay === 'true') {
                        video.play();
                    }
                } else {
                    if (!video.paused) {
                        video.pause();
                    }
                }
            });
        }

        window.addEventListener('scroll', handleVideoVisibility);
        window.addEventListener('resize', handleVideoVisibility);
        handleVideoVisibility();
    }

    // Подсветка активного пункта меню при прокрутке
    window.addEventListener('scroll', function () {
        const scrollPosition = window.scrollY;

        document.querySelectorAll('section[id]').forEach(section => {
            const sectionTop = section.offsetTop - 100;
            const sectionHeight = section.offsetHeight;
            const sectionId = section.getAttribute('id');

            if (scrollPosition >= sectionTop && scrollPosition < sectionTop + sectionHeight) {
                document.querySelectorAll('.navbar-nav .nav-link').forEach(link => {
                    link.classList.remove('active');
                    if (link.getAttribute('href') === '#' + sectionId) {
                        link.classList.add('active');
                    }
                });
            }
        });
    });
});
