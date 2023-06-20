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
    let $dinnerChoiceSelect = $("#dinner_choice");
    let $commentsInput = $("#comments");
    let $accommodations = $("#accommodations");
    let $guestNameInput = $("#guest_name");
    let $guestDinnerChoiceSelect = $("#guest_dinner_choice");

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

            $rsvpForm.find("input[name=guest_is_attending]:checked").prop("checked", false);
            $rsvpForm.find("input[name=guest_is_attending]:checked").prop("required", false);
        } else {
            $text.text("Remove Guest");

            $guestNameInput.prop("required", true);
            $guestNameInput.focus();

            $guestDinnerChoiceSelect.prop("required", true);

            $rsvpForm.find("input[name=guest_is_attending]:checked").prop("required", true);
        }
    });

    // Dinner Select required toggle based on attending radio
    $rsvpForm.find("input[name=is_attending]").change(function() {
        if ($rsvpForm.find("input[name=is_attending]:checked").val() === 'true') {
            $dinnerChoiceSelect.prop("required", true);
        } else {
            $dinnerChoiceSelect.prop("required", false);
        }
    });

    // Guest Dinner Select required toggle based on guest attending radio
    $rsvpForm.find("input[name=guest_is_attending]").change(function() {
        if ($rsvpForm.find("input[name=guest_is_attending]:checked").val() === 'true') {
            $guestDinnerChoiceSelect.prop("required", true);
        } else {
            $guestDinnerChoiceSelect.prop("required", false);
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
            is_attending: ($rsvpForm.find("input[name=is_attending]:checked").val() === 'true'),
            dinner_choice: parseInt($dinnerChoiceSelect.val()),
            accommodations: $accommodations.is(":checked"),
            comments: $commentsInput.val(),
        };

        if ($guestNameInput.prop("required") === true) {
            payload["guests"] = [
                {
                    name: $guestNameInput.val(),
                    dinner_choice: parseInt($guestDinnerChoiceSelect.val()),
                    is_attending: ($rsvpForm.find("input[name=guest_is_attending]:checked").val() === 'true'),
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
    $("#submitButton").prop("disabled", false);
}
onRecaptchaSubmit.name = "onRecaptchaSubmit"
