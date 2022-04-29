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
});
