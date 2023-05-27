$(document).ready(function() {
    let $countdownHolder = $("#countdown");

    // TODO: check timezone on this...
    let countDownDate = new Date("Oct 7, 2023 17:30:00").getTime();

    // Update the countdown every 1 second
    let x = setInterval(function() {
        let now = new Date().getTime();
        let diff = countDownDate - now;

        let days = Math.floor(diff / (1000 * 60 * 60 * 24));
        let hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
        let minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
        let seconds = Math.floor((diff % (1000 * 60)) / 1000);

        $countdownHolder.text(days + " days " + hours + " hours " + minutes + " minutes and " + seconds + " seconds");

        if (diff < 0) {
            clearInterval(x);
            $countdownHolder.css("visibility", "visible");
            $countdownHolder.text("Party Time!");
        }
    }, 1000);
});
