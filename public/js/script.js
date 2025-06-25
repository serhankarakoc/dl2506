// Sidebar ve backdrop
const sidebar = document.getElementById('sidebarMenu');
const sidebarToggle = document.getElementById('sidebarToggle');
const sidebarBackdrop = document.getElementById('sidebarBackdrop');
function closeSidebar() {
    sidebar.classList.remove('menu-open');
    sidebarBackdrop.classList.remove('active');
    document.body.classList.remove('sidebar-open');
}
function openSidebar() {
    sidebar.classList.add('menu-open');
    sidebarBackdrop.classList.add('active');
    document.body.classList.add('sidebar-open');
}
if (sidebarToggle && sidebar && sidebarBackdrop) {
    sidebarToggle.addEventListener('click', function (e) {
        e.stopPropagation();
        if (sidebar.classList.contains('menu-open')) {
            closeSidebar();
        } else {
            openSidebar();
        }
    });
    sidebarBackdrop.addEventListener('click', closeSidebar);
    // Sayfa içine tıklanınca sidebar kapansın
    document.addEventListener('click', function (e) {
        if (window.innerWidth <= 991 && sidebar.classList.contains('menu-open')) {
            // Sidebar veya toggle dışında bir yere tıklandıysa kapat
            if (!sidebar.contains(e.target) && e.target !== sidebarToggle) {
                closeSidebar();
            }
        }
    });
}
window.addEventListener('resize', function () {
    if (window.innerWidth > 991) closeSidebar();
});

// Profil dropdown menüsü
const profileBtn = document.getElementById('profileAvatarBtn');
const profileDropdown = document.getElementById('profileDropdown');
let profileDropdownOpen = false;
function closeProfileDropdown() {
    profileDropdown.classList.remove('show');
    profileDropdownOpen = false;
}
function openProfileDropdown() {
    profileDropdown.classList.add('show');
    profileDropdownOpen = true;
}
if (profileBtn && profileDropdown) {
    profileBtn.addEventListener('click', function (e) {
        e.stopPropagation();
        if (profileDropdownOpen) {
            closeProfileDropdown();
        } else {
            openProfileDropdown();
        }
    });
    document.addEventListener('click', function (e) {
        if (profileDropdownOpen && !profileDropdown.contains(e.target) && e.target !== profileBtn) {
            closeProfileDropdown();
        }
    });
}
window.addEventListener('resize', function () {
    if (window.innerWidth > 991) closeSidebar();
    closeProfileDropdown();
});

// Butonlara tıklandığında formu submit et ve butonu devre dışı bırak
document.addEventListener("DOMContentLoaded", function () {
    const buttons = document.querySelectorAll('button[type="submit"]');
    buttons.forEach((button) => {
        button.addEventListener("click", function () {
            button.disabled = true;
            button.innerHTML =
                '<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> Lütfen Bekleyiniz...';
            button.closest("form").submit();
        });
    });
});