document.getElementById("mssv").addEventListener("input", function() {
    var value = this.value;
    if (value.length > 8) {
        this.value = value.slice(0, 8);
    }
});

document.getElementById("confirm-password").addEventListener("blur", function() {
    var pass = document.getElementById("password").value;
    var confirm = this.value;
    var signupButton = document.getElementById("signup-button");
    if (pass !== confirm) {
        document.getElementById("error-message").innerHTML = "Passwords do not match!!!";
        document.getElementById("error-message").style.display = "block";
        signupButton.disabled = true;
    } else {
        document.getElementById("error-message").style.display = "none";
        signupButton.disabled = false;
    }
});
