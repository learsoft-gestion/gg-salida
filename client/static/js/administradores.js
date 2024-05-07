// import { prefijoURL } from './variables.js';
var prefijoURL

$(document).ready(function () {
    $('#menu').load('/static/menu.html');
});

fetch("/backend-url")
.then(res => res.json())
.then(data => {
    prefijoURL = data.prefijoURL;
    console.log(prefijoURL);   

    $.ajax({
        url: prefijoURL + '/convenios',
        method: 'GET',
        dataType: 'json',
        success: function (data) {
            if (data && data.length > 0) {
                data.forEach(convenio => {
                    const li = document.createElement("li");
                    li.value = convenio.id;
                    li.textContent = convenio.nombre;
                    $("#listaConvenios").append(li);
                });
            } else {
                console.log('No se recibieron datos del servidor.');
            }
        },
        error: function (error) {
            console.error('Error en la búsqueda:', error);
        }
    });
})
.catch(error => {
    console.error("Error al obtener la URL del backend: ", error)
})   