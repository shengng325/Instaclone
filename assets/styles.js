$( document ).ready(function() {
    loveBtn = $(".card-button.love");
    loveBtn.click(function(){
        btn = $(this).children()
        if (btn[0].style.display != "none") {
            console.log("1st")
            btn[0].style.display = "none"
            btn[1].style.display = "inline-block"
        } else {
            console.log("2nd")
            btn[0].style.display = "inline-block"
            btn[1].style.display = "none" 
        }
    });
});

