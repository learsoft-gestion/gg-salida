$(document).ready(function () {
    $('#menu').load('/static/menu.html');
});

$.ajax({
    url: '/convenios',
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