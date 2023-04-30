// ToDo: Make this multiple targeted files?
$(document).ready(function() {
    // Navbar
    let $burger = $(".navbar-burger");
    let active = $burger.attr("aria-expanded");
    $burger.click(function() {
        $burger.toggleClass("is-active");
        if (active === "true") {
            active = "false";
            $burger.attr("aria-expanded", active);
        } else {
            active = "true";
            $burger.attr("aria-expanded", active);
        }
        $(".navbar-menu").toggleClass("is-active");
    });

    // RSVP Form
    let $rsvpForm = $("#rsvpForm");
    let $addGuestButton = $("#addGuest");
    let $guestForm = $("#guestForm");

    // Form Fields
    let $nameInput = $("#name");
    let $emailInput = $("#email");
    let $attendingInput = $rsvpForm.find("input[name=is_attending]");
    let $dinnerChoiceSelect = $("#dinner_choice");
    let $commentsInput = $("#comments");
    let $guestNameInput = $("#guest_name");
    let $guestDinnerChoiceSelect = $("#guest_dinner_choice");

    // Form Guest Toggle
    $addGuestButton.click(function(e) {
        e.preventDefault();
        $guestForm.toggle();

        $addGuestButton.find("i").toggleClass("fa-user-plus fa-user-minus");

        let $text = $addGuestButton.find(".buttonText")
        if($guestForm.is(":hidden")) {
            $text.text("Add Guest");
            $guestNameInput.val("");
            $guestNameInput.prop("required", false);
            $guestDinnerChoiceSelect.val("");
            $guestDinnerChoiceSelect.prop("required", false);
        } else {
            $text.text("Remove Guest");
            $guestNameInput.prop("required", true);
            $guestNameInput.focus();
            $guestDinnerChoiceSelect.prop("required", true);
        }
    });

    // Form submit
    $rsvpForm.on("submit", function(e) {
        e.preventDefault();

        // ToDo: Prompt recaptcha

        let payload = {
            name: $nameInput.val(),
            email: $emailInput.val(),
            is_attending: ($attendingInput.val() === 'true'),
            dinner_choice: parseInt($dinnerChoiceSelect.val()),
            comments: $commentsInput.val(),
            guest_name: $guestNameInput.val(),
            guest_dinner_choice: parseInt($guestDinnerChoiceSelect.val()),
        }

        $.ajax({
            type: 'POST',
            url: '/api/rsvp',
            data: JSON.stringify(payload),
            contentType: "application/json; charset=utf-8",
            dataType: "json",
            success: function(data){
                // ToDo: Check for status created, Clear form, "Thanks for RSVPing" message, expect email?
                alert(data);
            },
            error: function(errMsg) {
                // ToDo: Post responseText message in an error banner with "contact mike if issues continue"
                alert(errMsg.responseText);
            },
        });

        return false;
    });
});
