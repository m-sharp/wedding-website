$(document).ready(function() {
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

    let $addGuest = $("#addGuest");
    let $guestForm = $("#guestForm");
    $addGuest.click(function (e) {
        e.preventDefault();
        $guestForm.toggle();

        $addGuest.find("i").toggleClass("fa-user-plus fa-user-minus");

        let $text = $addGuest.find(".buttonText")
        if($guestForm.is(":hidden")) {
            $text.text("Add Guest");
            $("#guestName").val("");
            $("#guestDinnerChoice").val("0");
        } else {
            $text.text("Remove Guest");
        }
    });
});
