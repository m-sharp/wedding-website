$(document).ready(function() {
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
    let $accommodations = $("#accommodations");
    let $guestNameInput = $("#guest_name");
    let $guestDinnerChoiceSelect = $("#guest_dinner_choice");
    let $guestAttendingInput = $rsvpForm.find("input[name=guest_is_attending]");

    // Checkbox text accessibility
    $("#accommodationSpan").click(function(e) {
        $accommodations.click();
    })

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

            $guestAttendingInput.prop("required", true)
        }
    });

    // Form submit
    $rsvpForm.on("submit", function(e) {
        e.preventDefault();

        let payload = {
            name: $nameInput.val(),
            email: $emailInput.val(),
            is_attending: ($attendingInput.val() === 'true'),
            dinner_choice: parseInt($dinnerChoiceSelect.val()),
            accommodations: $accommodations.is(":checked"),
            comments: $commentsInput.val(),
        }

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
