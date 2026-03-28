// Read the CSRF token from the csrf_token cookie
function getCsrfToken() {
    const match = document.cookie.match(/(^|;\s*)csrf_token=([^;]*)/);
    return match ? decodeURIComponent(match[2]) : '';
}
