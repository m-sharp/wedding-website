$(document).ready(function() {
    let $countdownHolder = $("#countdown");
    let countDownDate = new Date("Oct 7, 2023 17:30:00").getTime();

    let pluralize = function(noun, val) {
        if (val === 1) {
            return " " + val + " " + noun + " ";
        }
        return " " + val + " " + noun + "s ";
    };

    let setCountDown = function() {
        let now = new Date().getTime();
        let diff = countDownDate - now;
        if (diff < 0) {
            $countdownHolder.text("Party Time!");
            clearInterval(x);
            return;
        }

        let days = Math.floor(diff / (1000 * 60 * 60 * 24));
        let hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
        let minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
        let seconds = Math.floor((diff % (1000 * 60)) / 1000);

        let toSet = "";
        if (days > 0) {
            toSet += pluralize("day", days);
        }
        if (hours > 0) {
            toSet += pluralize("hour", hours);
        }
        if (minutes >0) {
            toSet += pluralize("minute", minutes);
        }
        toSet += pluralize("second", seconds);
        $countdownHolder.text(toSet);
    };

    // Run setCountDown immediately
    setCountDown();
    // Update the countdown every 1 second
    let x = setInterval(setCountDown, 1000);
});
