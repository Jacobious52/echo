window.onload = function() {
    var user = document.getElementById("user");

    document.getElementById("login").onsubmit = function() {
        if (!user.value) {
            return false;
        }

        if (user.value.trim().length < 3) {
            alert("username must be longer than 3 valid characters");
            return false;
        }

        window.location.href += "/chat.html#" + user.value;

        return false;
    };
}
