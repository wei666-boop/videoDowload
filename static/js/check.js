const type=document.getElementById("select-type")
const format=document.getElementById("text-format");
const subtitle=document.getElementById("text-subtitle")
const thumbnail=document.getElementById("text-thumbnail")

function updateUI(){
    const isAudio=type.value==="audio"
    const isAuto=format.value==="auto"

    if(isAuto===true){
        subtitle.disabled=true
        subtitle.checked=false
        thumbnail.disabled=true
        thumbnail.checked=false
    }
    else if(isAuto===false){
        if(isAudio===true){
            subtitle.disabled=true
            subtitle.checked=false
            thumbnail.disabled=true
            thumbnail.checked=false
        }
        else if(isAudio===false){
            subtitle.disabled=false
            thumbnail.disabled=false
        }
    }
}

type.addEventListener("change",updateUI)
format.addEventListener("change",updateUI)


updateUI()
