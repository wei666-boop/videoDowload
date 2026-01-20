async function loadRecords(){
    const response=await fetch("http://localhost:5443/center/record")
    const records=await response.json()

    const tbody=document.getElementById("record-table-body")
    tbody.innerText=""

    records.forEach(record=>{
        const tr=document.createElement("tr")
        tr.innerHTML=`
        <td>${record.id}</td>
        <td>${record.url}</td>
        <td>${record.time}</td>
            `
        tbody.appendChild(tr)
    })
    
}

function clear (){
    const td=document.getElementById("record-table-body")
    td.innerHTML=""
    fetch("http://localhost:5443/center/clear")
        .then((res)=>{
            if(!res.ok){
                alert("清空失败")
            }
        })
        .catch(err=>{
            console.error(err)
        })
}