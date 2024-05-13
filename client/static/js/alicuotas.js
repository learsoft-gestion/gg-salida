var prefijoURL

$('#menu').load('/static/menu.html', function () {
    $('#titulo').append('Alícuotas');
});

fetch("/backend-url")
    .then(res => res.json())
    .then(data => {
        prefijoURL = data.prefijoURL;
        console.log(prefijoURL);

        // Select de Sindicato
        $.ajax({
            url: prefijoURL + '/convenios',
            method: 'GET',
            dataType: 'json',
            success: function (data) {
                if (data && data.length > 0) {
                    data.forEach(convenio => {
                        const option = document.createElement("option");
                        option.value = convenio.id;
                        option.textContent = convenio.nombre;
                        $("#sindicato").append(option);
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
    });

// Llamado a alícuotas
$('#sindicato').change(function () {
    $.ajax({
        url: prefijoURL + `/alicuotas/${$('#sindicato').val()}`,
        method: 'GET',
        dataType: 'json',
        success: function (data) {
            mostrarAlicuotas(data);
            $('#valores').hide();
            $('#valoresData').empty();
        },
        error: function (error) {
            console.error('Error en la búsqueda:', error);
            Swal.fire({
                title: "Ocurrió un error",
                text: error.mensaje,
                icon: "error"
            });
        }
    });
});

var mostrarAlicuotas = function (data) {
    $('#alicuotasData').empty();

    $.each(data, function (index, alicuota) {
        const group = $(`<div class="form-group" id="${alicuota.IdAlicuota}">`);
        const control = $(`<div class="form-control">`);
        const nombre = $(`<input type="text" disabled value="${alicuota.Nombre}" id="nombre-${alicuota.IdAlicuota}" class="nombre">`);
        const descripcion = $(`<input type="text" disabled value="${alicuota.Descripcion}" id="descripcion-${alicuota.IdAlicuota}" class="descripcion">`);
        const acciones = $(`<div class="actions">`);
        const btnEdit = $(`<button type="button" class="btn btn-sm" id="editAli-${alicuota.IdAlicuota}" title="Editar">`);
        const spanEdit = $(`<span class="material-symbols-outlined">edit</span>`);
        const btnValue = $(`<button type="button" class="btn btn-sm" id="valuesAli-${alicuota.IdAlicuota}" title="Valores">`);
        const spanValue = $(`<span class="material-symbols-outlined">double_arrow</span>`);
        const btnSave = $(`<button type="button" class="btn btn-sm d-none" id="saveAli-${alicuota.IdAlicuota}" title="Guardar">`);
        const spanSave = $(`<span class="material-symbols-outlined">done_outline</span>`);
        const btnDelete = $(`<button type="button" class="btn btn-sm" id="deleteAli-${alicuota.IdAlicuota}" title="Borrar">`);
        const spanDelete = $(`<span class="material-symbols-outlined">delete</span>`);

        btnDelete.append(spanDelete);
        btnSave.append(spanSave);
        btnValue.append(spanValue);
        btnEdit.append(spanEdit);
        acciones.append(btnEdit);
        acciones.append(btnSave);
        acciones.append(btnValue);
        acciones.append(btnDelete);
        control.append(nombre);
        control.append(descripcion);
        control.append(acciones);
        group.append(control);

        $('#alicuotasData').append(group);
        botonEditAli(alicuota.IdAlicuota);
        botonValuesAli(alicuota.IdAlicuota);
        botonSaveAli(alicuota.IdAlicuota);
        botonDeleteAli(alicuota.IdAlicuota);
    });

    $('#alicuotas').show();
}

$('#btnAddAlicuota').click(function () {
    if ($('#saveAli-').is(':visible')) {
        return;
    }

    const group = $(`<div class="form-group" id="newAli">`);
    const control = $(`<div class="form-control">`);
    const nombre = $(`<input type="text" placeholder="Nombre" id="nombre-" class="nombre">`);
    const descripcion = $(`<input type="text" placeholder="Descripción" id="descripcion-" class="descripcion">`);
    const acciones = $(`<div class="actions">`);
    const btnEdit = $(`<button type="button" class="btn btn-sm d-none" id="editAli-" title="Editar">`);
    const spanEdit = $(`<span class="material-symbols-outlined">edit</span>`);
    const btnValue = $(`<button type="button" class="btn btn-sm d-none" id="valuesAli-" title="Valores">`);
    const spanValue = $(`<span class="material-symbols-outlined">double_arrow</span>`);
    const btnSave = $(`<button type="button" class="btn btn-sm" id="saveAli-" title="Guardar">`);
    const spanSave = $(`<span class="material-symbols-outlined">done_outline</span>`);
    const btnDelete = $(`<button type="button" class="btn btn-sm d-none" id="deleteAli-" title="Borrar">`);
    const spanDelete = $(`<span class="material-symbols-outlined">delete</span>`);
    const btnCancel = $(`<button type="button" class="btn btn-sm" id="cancelAli" title="Cancelar">`);
    const spanCancel = $(`<span class="material-symbols-outlined">close</span>`);

    btnCancel.append(spanCancel);
    btnDelete.append(spanDelete);
    btnSave.append(spanSave);
    btnValue.append(spanValue);
    btnEdit.append(spanEdit);
    acciones.append(btnEdit);
    acciones.append(btnValue);
    acciones.append(btnSave);
    acciones.append(btnDelete);
    acciones.append(btnCancel);
    control.append(nombre);
    control.append(descripcion);
    control.append(acciones);
    group.append(control);

    $('#alicuotasData').append(group);

    botonCancelAli();
    botonCreateAli();
});

var botonCancelAli = function () {
    $('#cancelAli').click(function () {
        $('#newAli').remove();
    });
}

var botonCreateAli = function () {
    $('#saveAli-').click(function () {
        var json = {
            idConvenio: $(`#sindicato`).val(),
            nombre: $(`#nombre-`).val(),
            descripcion: $(`#descripcion-`).val()
        }

        $.ajax({
            url: prefijoURL + `/alicuotas`,
            method: 'POST',
            dataType: 'json',
            data: JSON.stringify(json),
            success: function (data) {
                if (data) {
                    $(`#saveAli-`).attr('id', `saveAli-${data.id}`);
                    $(`#editAli-`).attr('id', `editAli-${data.id}`);
                    $(`#deleteAli-`).attr('id', `deleteAli-${data.id}`);
                    $(`#nombre-`).attr('id', `nombre-${data.id}`);
                    $(`#descripcion-`).attr('id', `descripcion-${data.id}`);
                    $(`#saveAli-${data.id}`).hide();
                    $(`#editAli-${data.id}`).removeClass('d-none');
                    $(`#daleteAli-${data.id}`).removeClass('d-none');
                    $(`#nombre-${data.id}`).attr('disabled', true);
                    $(`#descripcion-${data.id}`).attr('disabled', true);
                    botonEditAli(data.id);
                    botonValuesAli(data.id);
                    botonSaveAli(data.id);
                    botonDeleteAli(data.id);
                } else {
                    Swal.fire("No se pudo editar la alícuota");
                }
            },
            error: function (error) {
                console.error('Error en la búsqueda:', error);
                Swal.fire({
                    title: "Ocurrió un error",
                    text: error.mensaje,
                    icon: "error"
                });
            }
        });
    });
}

var botonEditAli = function (id) {
    $(`#editAli-${id}`).click(function () {
        $(`#nombre-${id}`).removeAttr('disabled');
        $(`#descripcion-${id}`).removeAttr('disabled');

        $(`#editAli-${id}`).hide();
        $(`#saveAli-${id}`).removeClass('d-none');
        $(`#saveAli-${id}`).show();
    });
}

var botonValuesAli = function (id) {
    $(`#valuesAli-${id}`).click(function () {
        $.ajax({
            url: prefijoURL + `/valoresAlicuotas/${id}`,
            method: 'GET',
            dataType: 'json',
            success: function (data) {
                $(`.form-control`).removeClass('bg-amarillo');
                $(`#${id} .form-control`).addClass('bg-amarillo');
                mostrarValoresAlicuotas(data, id);
            },
            error: function (error) {
                console.error('Error en la búsqueda:', error);
                Swal.fire({
                    title: "Ocurrió un error",
                    text: error.mensaje,
                    icon: "error"
                });
            }
        });
    });
}

var botonSaveAli = function (id) {
    $(`#saveAli-${id}`).click(function () {
        var json = {
            idAlicuota: id,
            nombre: $(`#nombre-${id}`).val(),
            descripcion: $(`#descripcion-${id}`).val()
        }

        $.ajax({
            url: prefijoURL + `/alicuotas`,
            method: 'PATCH',
            dataType: 'json',
            data: JSON.stringify(json),
            success: function (data) {
                if (data) {
                    $(`#saveAli-${id}`).hide();
                    $(`#editAli-${id}`).show();
                    $(`#nombre-${id}`).attr('disabled', true);
                    $(`#descripcion-${id}`).attr('disabled', true);
                } else {
                    Swal.fire("No se pudo editar la alícuota");
                }
            },
            error: function (error) {
                console.error('Error en la búsqueda:', error);
                Swal.fire({
                    title: "Ocurrió un error",
                    text: error.mensaje,
                    icon: "error"
                });
            }
        });
    });
}

var botonDeleteAli = function (id) {
    $(`#deleteAli-${id}`).click(function () {
        Swal.fire({
            title: "¿Quiere borrar la alícuota?",
            text: "No podrá revertir esta acción",
            icon: "warning",
            showCancelButton: true,
            confirmButtonColor: "#3085d6",
            cancelButtonColor: "#d33",
            confirmButtonText: "Sí, borralo",
            cancelButtonText: "Cancelar"
        }).then((result) => {
            if (result.isConfirmed) {
                $.ajax({
                    url: prefijoURL + `/alicuotas`,
                    method: 'DELETE',
                    dataType: 'json',
                    data: JSON.stringify({ idAlicuota: id }),
                    success: function (data) {
                        if (data) {
                            $(`#sindicato`).trigger('change');
                            Swal.fire({
                                title: "¡Borrado!",
                                text: "La alícuota se ha borrado.",
                                icon: "success"
                            });
                        } else {
                            Swal.fire("No se pudo borrar la alícuota");
                        }
                    },
                    error: function (error) {
                        console.error('Error en la búsqueda:', error);
                        Swal.fire({
                            title: "Ocurrió un error",
                            text: error.mensaje,
                            icon: "error"
                        });
                    }
                });
            }
        });
    });
}

var mostrarValoresAlicuotas = function (data, idAlicuota) {
    $('#valoresData').empty();
    $('#valoresData').val(idAlicuota);

    $.each(data, function (index, valor) {
        const group = $(`<div class="form-group" id="${valor.IdValoresAlicuota}">`);
        const control = $(`<div class="form-control">`);
        const labelPeriodo = $(`<label for="periodo-${valor.IdValoresAlicuota}">Período</label>`);
        const periodo = $(`<input type="text" disabled value="${valor.VigenciaDesde}" id="periodo-${valor.IdValoresAlicuota}" class="periodo">`);
        const labelValor = $(`<label for="valor-${valor.IdValoresAlicuota}">Valor</label>`);
        const valorAli = $(`<input type="text" disabled value="${valor.Valor}" id="valor-${valor.IdValoresAlicuota}" class="valor">`);
        const acciones = $(`<div class="actions">`);
        const btnEdit = $(`<button type="button" class="btn btn-sm" id="editVal-${valor.IdValoresAlicuota}" title="Editar">`);
        const spanEdit = $(`<span class="material-symbols-outlined">edit</span>`);
        const btnSave = $(`<button type="button" class="btn btn-sm d-none" id="saveVal-${valor.IdValoresAlicuota}" title="Guardar">`);
        const spanSave = $(`<span class="material-symbols-outlined">done_outline</span>`);
        const btnDelete = $(`<button type="button" class="btn btn-sm" id="deleteVal-${valor.IdValoresAlicuota}" title="Borrar">`);
        const spanDelete = $(`<span class="material-symbols-outlined">delete</span>`);

        btnDelete.append(spanDelete);
        btnSave.append(spanSave);
        btnEdit.append(spanEdit);
        acciones.append(btnEdit);
        acciones.append(btnSave);
        acciones.append(btnDelete);
        control.append(labelPeriodo);
        control.append(periodo);
        control.append(labelValor);
        control.append(valorAli);
        control.append(acciones);
        group.append(control);

        $('#valoresData').append(group);
        botonEditValor(valor.IdValoresAlicuota);
        botonSaveValor(valor.IdValoresAlicuota);
        botonDeleteValor(valor.IdValoresAlicuota);
    });

    $(".periodo").datepicker({
        autoclose: true,
        minViewMode: 1,
        format: 'yyyymm',
        language: "es"
    });

    $('#valores').show();
}

var botonEditValor = function (id) {
    $(`#editVal-${id}`).click(function () {
        $(`#periodo-${id}`).removeAttr('disabled');
        $(`#valor-${id}`).removeAttr('disabled');

        $(`#editVal-${id}`).hide();
        $(`#saveVal-${id}`).removeClass('d-none');
        $(`#saveVal-${id}`).show();
    });
}

var botonSaveValor = function (id) {
    $(`#saveVal-${id}`).click(function () {
        var json = {
            idValoresAlicuota: String(id),
            vigenciaDesde: $(`#periodo-${id}`).val(),
            valor: $(`#valor-${id}`).val()
        }

        $.ajax({
            url: prefijoURL + `/valoresAlicuotas`,
            method: 'PATCH',
            dataType: 'json',
            data: JSON.stringify(json),
            success: function (data) {
                if (data) {
                    $(`#saveVal-${id}`).hide();
                    $(`#editVal-${id}`).show();
                    $(`#periodo-${id}`).attr('disabled', true);
                    $(`#valor-${id}`).attr('disabled', true);
                } else {
                    Swal.fire("No se pudo editar el valor");
                }
            },
            error: function (error) {
                console.error('Error en la búsqueda:', error);
                Swal.fire({
                    title: "Ocurrió un error",
                    text: error.mensaje,
                    icon: "error"
                });
            }
        });
    });
}

var botonDeleteValor = function (id) {
    $(`#deleteVal-${id}`).click(function () {
        Swal.fire({
            title: "¿Quiere borrar el valor?",
            text: "No podrá revertir esta acción",
            icon: "warning",
            showCancelButton: true,
            confirmButtonColor: "#3085d6",
            cancelButtonColor: "#d33",
            confirmButtonText: "Sí, borralo",
            cancelButtonText: "Cancelar"
        }).then((result) => {
            if (result.isConfirmed) {
                $.ajax({
                    url: prefijoURL + `/valoresAlicuotas`,
                    method: 'DELETE',
                    dataType: 'json',
                    data: JSON.stringify({ idValoresAlicuota: id }),
                    success: function (data) {
                        if (data) {
                            $(`#${id}`).remove();
                            Swal.fire({
                                title: "¡Borrado!",
                                text: "El valor se ha borrado.",
                                icon: "success"
                            });
                        } else {
                            Swal.fire("No se pudo borrar el valor");
                        }
                    },
                    error: function (error) {
                        console.error('Error en la búsqueda:', error);
                        Swal.fire({
                            title: "Ocurrió un error",
                            text: error.mensaje,
                            icon: "error"
                        });
                    }
                });
            }
        });
    });
}

$('#btnAddValor').click(function () {
    if ($('#saveVal-').is(':visible')) {
        return;
    }

    const group = $(`<div class="form-group" id="newVal">`);
    const control = $(`<div class="form-control">`);
    const labelPeriodo = $(`<label for="periodo-">Período</label>`);
    const periodo = $(`<input type="text" value="" id="periodo-" class="periodo">`);
    const labelValor = $(`<label for="valor-">Valor</label>`);
    const valorAli = $(`<input type="text" value="" id="valor-" class="valor">`);
    const acciones = $(`<div class="actions">`);
    const btnEdit = $(`<button type="button" class="btn btn-sm d-none" id="editVal-" title="Editar">`);
    const spanEdit = $(`<span class="material-symbols-outlined">edit</span>`);
    const btnSave = $(`<button type="button" class="btn btn-sm" id="saveVal-" title="Guardar">`);
    const spanSave = $(`<span class="material-symbols-outlined">done_outline</span>`);
    const btnDelete = $(`<button type="button" class="btn btn-sm d-none" id="deleteVal-" title="Borrar">`);
    const spanDelete = $(`<span class="material-symbols-outlined">delete</span>`);
    const btnCancel = $(`<button type="button" class="btn btn-sm" id="cancelVal" title="Cancelar">`);
    const spanCancel = $(`<span class="material-symbols-outlined">close</span>`);

    btnCancel.append(spanCancel);
    btnDelete.append(spanDelete);
    btnSave.append(spanSave);
    btnEdit.append(spanEdit);
    acciones.append(btnEdit);
    acciones.append(btnSave);
    acciones.append(btnDelete);
    acciones.append(btnCancel);
    control.append(labelPeriodo);
    control.append(periodo);
    control.append(labelValor);
    control.append(valorAli);
    control.append(acciones);
    group.append(control);

    $('#valoresData').append(group);

    $("#periodo-").datepicker({
        autoclose: true,
        minViewMode: 1,
        format: 'yyyymm',
        language: "es"
    });
    botonCancelVal();
    botonCreateVal();
});

var botonCancelVal = function () {
    $('#cancelVal').click(function () {
        $('#newVal').remove();
    });
}

var botonCreateVal = function () {
    $('#saveVal-').click(function () {
        var json = {
            idAlicuota: $(`#valoresData`).val(),
            vigenciaDesde: $(`#periodo-`).val(),
            valor: $(`#valor-`).val()
        }

        $.ajax({
            url: prefijoURL + `/valoresAlicuotas`,
            method: 'POST',
            dataType: 'json',
            data: JSON.stringify(json),
            success: function (data) {
                if (data) {
                    $(`#newVal-`).attr('id', `${data.id}`);
                    $(`#saveVal-`).attr('id', `saveVal-${data.id}`);
                    $(`#editVal-`).attr('id', `editVal-${data.id}`);
                    $(`#deleteVal-`).attr('id', `deleteVal-${data.id}`);
                    $(`#periodo-`).attr('id', `periodo-${data.id}`);
                    $(`#valor-`).attr('id', `valor-${data.id}`);
                    $(`#saveVal-${data.id}`).hide();
                    $(`#editVal-${data.id}`).removeClass('d-none');
                    $(`#deleteVal-${data.id}`).removeClass('d-none');
                    $(`#periodo-${data.id}`).attr('disabled', true);
                    $(`#valor-${data.id}`).attr('disabled', true);
                    botonEditValor(data.id);
                    botonSaveValor(data.id);
                    botonDeleteValor(data.id);
                } else {
                    Swal.fire("No se pudo editar el valor");
                }
            },
            error: function (error) {
                console.error('Error en la búsqueda:', error);
                Swal.fire({
                    title: "Ocurrió un error",
                    text: error.mensaje,
                    icon: "error"
                });
            }
        });
    });
}