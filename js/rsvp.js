$(document).ready(function() {
    // RSVP Form
    let $rsvpForm = $("#rsvpForm");

    let $errorBanner = $("#errorBanner");
    let $errorText = $("#errorText");

    let $successText = $("#successText");

    let $addGuestButton = $("#addGuest");
    let $guestForm = $("#guestForm");

    // Form Fields
    let $nameInput = $("#name");
    let $emailInput = $("#email");
    let $attendingInput = $rsvpForm.find("input[name=is_attending]:checked");
    let $dinnerChoiceSelect = $("#dinner_choice");
    let $commentsInput = $("#comments");
    let $accommodations = $("#accommodations");
    let $guestNameInput = $("#guest_name");
    let $guestDinnerChoiceSelect = $("#guest_dinner_choice");
    let $guestAttendingInput = $rsvpForm.find("input[name=guest_is_attending]:checked");

    // Security
    let $recaptcha = $("#recaptcha");
    let $csrfToken = $("input[name='gorilla.csrf.Token']");

    // Checkbox text accessibility
    $("#accommodationSpan").click(function(e) {
        $accommodations.click();
    });

    // Error close
    $("#errorClose").click(function(e) {
        $errorBanner.hide();
    });

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

            $guestAttendingInput.prop("checked", false);
            $guestAttendingInput.prop("required", false);
        } else {
            $text.text("Remove Guest");

            $guestNameInput.prop("required", true);
            $guestNameInput.focus();

            $guestDinnerChoiceSelect.prop("required", true);

            $guestAttendingInput.prop("required", true);
        }
    });

    // Form submit
    $rsvpForm.on("submit", function(e) {
        e.preventDefault();

        let responseToken = grecaptcha.getResponse();
        if (responseToken === "") {
            $recaptcha.addClass("recaptchaNeeded");
            return;
        } else {
            $recaptcha.removeClass("recaptchaNeeded");
        }

        let payload = {
            name: $nameInput.val(),
            email: $emailInput.val(),
            is_attending: ($attendingInput.val() === 'true'),
            dinner_choice: parseInt($dinnerChoiceSelect.val()),
            accommodations: $accommodations.is(":checked"),
            comments: $commentsInput.val(),
        };

        if ($guestNameInput.prop("required") === true) {
            payload["guests"] = [
                {
                    name: $guestNameInput.val(),
                    dinner_choice: parseInt($guestDinnerChoiceSelect.val()),
                    is_attending: ($guestAttendingInput.val() === 'true'),
                }
            ];
        }

        $.ajax({
            type: 'POST',
            url: '/api/rsvp?' + $.param({token: responseToken}),
            data: JSON.stringify(payload),
            contentType: "application/json; charset=utf-8",
            headers: {
                "X-CSRF-Token": $csrfToken.val(),
            },
            success: function(data, textStatus){
                $errorBanner.hide();

                $rsvpForm[0].reset();
                $rsvpForm.css("visibility", "hidden");

                $successText.show();
                $successText.addClass("is-flex");

                grecaptcha.reset();
            },
            error: function(errMsg) {
                $errorText.text("Error submitting form: " + errMsg.responseText.trim() + ".");
                $errorBanner.show();
            },
            complete: function () {
                window.scrollTo(0, 0);
            },
        });

        return false;
    });
});

// Recaptcha doesn't like `let`, so `var` it is
var onRecaptchaSubmit = function() {
    $("#submitButton").prop('disabled', false);
}
onRecaptchaSubmit.name = "onRecaptchaSubmit"
