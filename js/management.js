$(document).ready(function() {
    let $errorBanner = $("#errorBanner");

    // Error close
    $("#errorClose").click(function(e) {
        $errorBanner.hide();
    });
});

let remove = function(id) {
    let $errorBanner = $("#errorBanner");
    let $errorText = $("#errorText");
    let $csrfToken = $("input[name='gorilla.csrf.Token']");

    let $targetRow = $("#row-" + id);
    let $spinner = $targetRow.find(".spinner");
    let $delete = $targetRow.find(".trash");

    if (confirm("Are you sure you want to delete this RSVP?") !== true) {
        return;
    }

    $delete.hide();
    $spinner.show();

    $.ajax({
        type: 'DELETE',
        url: '/api/rsvp/' + id,
        headers: {
            "X-CSRF-Token": $csrfToken.val(),
        },
        success: function(data, textStatus){
            $errorBanner.hide();
            $targetRow.hide();
        },
        error: function(errMsg) {
            $errorText.text("Error deleting record: " + errMsg.responseText.trim() + ".");
            $errorBanner.show();

            $spinner.hide();
            $delete.show();
        },
        complete: function () {
            window.scrollTo(0, 0);
        },
    });

    return false;
};
