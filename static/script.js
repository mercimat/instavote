var elemA = document.getElementById("a")
var elemB = document.getElementById("b")
var url = window.location.origin + "/api/results"

function get() {
    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);
    xhr.onreadystatechange = callback;
    xhr.send();

    function callback() {
        if (xhr.readyState == 4 && xhr.status == 200) {
            var response = JSON.parse(xhr.responseText);
            for (i in response) {
                if (i == "a") {
                    elemA.innerHTML = response[i]
                } else {
                    elemB.innerHTML = response[i]
                }
            }
        }
    }
}

setInterval(get, 1000);
