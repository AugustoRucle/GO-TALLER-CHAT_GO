$(document).ready(function (){
    var user_name
    var final_conexion

    $("#form_registro").on("submit", function(e){
        e.preventDefault(); //No se ejecuta el submit
        user_name= $("#user_name").val()

        $.ajax({
            type:"POST",
            url:"http://localhost:8000/validate",
            data:{
                "user_name": user_name
            },
            success: function(data){
                result(data)
            }
        })
    })

    function result(data){
        console.log(data)
        obj = JSON.parse(data) // --> ¿Que es parse?
        if(obj.IsValid){
            create_conexion()
        }else{
            console.log("Intenlo de nuevo :3")
        }
    }

    function create_conexion(){
        $("#registro").hide()
        $("#container_chat").show()
        var conexion = new WebSocket("ws://localhost:8000/chat/" + user_name)
        final_conexion = conexion
        conexion.onopen = function (response){
            conexion.onmessage = function(response){
                console.log(response.data)
                val = $("#chat_area").val()
                $("#chat_area").val(val + "\n" + response.data)
            }
        }
    }

    $("#form_message").on("submit", function(e){
        e.preventDefault();
        mensaje = $("#msg").val()
        final_conexion.send(mensaje)
        $("#msg").val("")
    })
})
