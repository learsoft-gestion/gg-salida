$(document).ready(function () {
    // var fetchedData;
    const selectConv = document.getElementById("conv");
    const selectEmp = document.getElementById("emp");
    const selectProc = document.getElementById("select");

    document.addEventListener("DOMContentLoaded", function () {
        // Realizar la solicitud de red (fetch) al servidor
        fetch("/convenios")
            .then(response => {
                // response.json().then(console.log(...data))
                if (!response.ok) {
                    throw new Error('La solicitud falló: ' + response.status);
                }
                return response.json();
            })
            .then(data => {
                // fetchedData = data
                // Agregar nuevas opciones al select
                data.forEach(convenio => {
                    const option = document.createElement("option");
                    option.value = convenio.id; // Asegúrate de que proceso tenga una propiedad "id"
                    option.textContent = convenio.nombre; // Asegúrate de que proceso tenga una propiedad "nombre"
                    selectConv.appendChild(option);
                });
            })
            .catch(error => {
                console.error("Error al obtener los procesos:", error);
            });
    });

    var selectedConv
    document.getElementById("conv").addEventListener("change", function () {
        selectedConv = this.options[this.selectedIndex];
        // console.log("Opción seleccionada:", selectedOption.textContent);

        // Limpiar opciones existentes
        const hijos = selectEmp.querySelectorAll(":not(:first-child)");
        hijos.forEach(hijo => {
            hijo.remove();
        });

        // Realizar la solicitud de red (fetch) al servidor
        fetch(`/empresas/${selectConv.value}`)
            .then(response => {
                // response.json().then(console.log(...data))
                if (!response.ok) {
                    throw new Error('La solicitud falló: ' + response.status);
                }
                return response.json();
            })
            .then(data => {
                // fetchedData = data
                // Agregar nuevas opciones al select
                data.forEach(empresa => {
                    const option = document.createElement("option");
                    option.value = empresa.id; // Asegúrate de que proceso tenga una propiedad "id"
                    option.textContent = empresa.nombre; // Asegúrate de que proceso tenga una propiedad "nombre"
                    selectEmp.appendChild(option);
                });
            })
            .catch(error => {
                console.error("Error al obtener los procesos:", error);
            });
    });

    var selectedEmp
    document.getElementById("emp").addEventListener("change", function () {
        selectedEmp = this.options[this.selectedIndex];
        // console.log("Opción seleccionada:", selectedEmp.value);

        // Limpiar opciones existentes
        const hijos = selectProc.querySelectorAll(":not(:first-child)");
        hijos.forEach(hijo => {
            hijo.remove();
        });

        // Realizar la solicitud de red (fetch) al servidor
        fetch(`/procesos/${selectConv.value}/${selectEmp.value}`)
            .then(response => {
                // response.json().then(console.log(...data))
                if (!response.ok) {
                    throw new Error('La solicitud falló: ' + response.status);
                }
                return response.json();
            })
            .then(data => {
                // Agregar nuevas opciones al select
                data.forEach(proceso => {
                    const option = document.createElement("option");
                    option.value = proceso.id; // Asegúrate de que proceso tenga una propiedad "id"
                    option.textContent = proceso.nombre; // Asegúrate de que proceso tenga una propiedad "nombre"
                    selectProc.appendChild(option);
                });
            })
            .catch(error => {
                console.error("Error al obtener los procesos:", error);
            });
    });

    var selectedProc
    selectProc.addEventListener("change", function () {
        selectedProc = this.options[this.selectedIndex];
        var li = document.createElement("li");
        li.appendChild(document.createTextNode(selectedProc.textContent))
        li.id = selectedProc.value

        var btnEliminar = document.createElement("button")
        btnEliminar.appendChild(document.createTextNode("X"))
        btnEliminar.addEventListener("click", function () {
            li.remove()
        });

        li.appendChild(btnEliminar)

        document.getElementById("listado").appendChild(li)
    });

    periodoCheck.addEventListener("click", () => {
        fecha2.style.display = 'block'
        inicial.style.display = 'block'
        fechaUnica.style.display = 'none'
    });
    fechaCheck.addEventListener("click", () => {
        fecha2.style.display = 'none'
        inicial.style.display = 'none'
        fechaUnica.style.display = 'block'
    });


    // Función para mostrar la confirmación
    function mostrarConfirmacion() {
        const selectValue = document.getElementById("select").value
        const fecha = document.getElementById("fecha").value
        const fechaFormat = fecha.replaceAll("-", "")
        const fecha1 = document.getElementById("fecha1").value
        const fechaFormat2 = fecha1.replaceAll("-", "")
        Swal.fire({
            title: 'El modelo fue procesado anteriormente, desea continuar?',
            showDenyButton: true,
            showCancelButton: false,
            confirmButtonText: 'Si',
            denyButtonText: 'No',
            customClass: {
                actions: 'my-actions',
                cancelButton: 'order-1 right-gap',
                confirmButton: 'order-2',
                denyButton: 'order-3',
            },
        }).then((result) => {
            if (result.isConfirmed) {
                // Swal.fire('Saved!', '', 'success')
                fetch("/force", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application-json"
                    },
                    body: JSON.stringify({
                        id: Number(selectValue),
                        fecha: fechaFormat,
                        fecha2: fechaFormat2,
                        forzado: true
                    })
                })
                    .then(resp => resp.json())
                    .then(data => {

                        if (data.archivos_salida == null) {
                            mensaje.textContent = data.mensaje
                            mensaje.style.display = 'block'
                            btnEnviar.disabled = false
                            loader.style.display = 'none'
                        } else {
                            data.archivos_salida.forEach(archivo => {
                                var span = document.createElement("span")
                                span.textContent = "Archivo guardado en: " + archivo
                                mensaje.appendChild(span)
                            })
                            // mensaje.textContent = "El archivo se guardó en: " + data["archivo_salida"]
                            ensaje.style.display = 'block'
                            tnEnviar.disabled = false              // Habilito el boton
                            oader.style.display = 'none'           // Saco loader
                        }

                    })
            } else if (result.isDenied) {
                // Swal.fire('Selecciona el proceso que desee ejecutar', '', 'info')
                mensaje.style.display = 'block'
                btnEnviar.disabled = false
                loader.style.display = 'none'
            }
        })
    }

    var seleccionados = []
    const ul = document.getElementById("listado")
    const btnEnviar = document.getElementById("btnEnviar")

    btnEnviar.addEventListener("click", () => {
        var elementos = ul.getElementsByTagName("li")
        for (var i = 0; i < elementos.length; i++) {
            var texto = elementos[i].id
            seleccionados.push(Number(texto))
        }

        btnEnviar.disabled = true
        loader.style.display = 'block'

        // const selectValue = document.getElementById("select").value
        const fecha = document.getElementById("fecha").value
        const fechaFormat = fecha.replaceAll("-", "")
        const fecha1 = document.getElementById("fecha1").value
        const fechaFormat2 = fecha1.replaceAll("-", "")

        fetch("/send", {
            method: "POST",
            headers: {
                "Content-Type": "application-json"
            },
            body: JSON.stringify({
                ids: seleccionados,
                fecha: fechaFormat,
                fecha2: fechaFormat2,
            })
        })
            .then(resp => resp.json())
            .then(data => {

                // if (data.procesado) {
                //   mostrarConfirmacion()
                // } else 
                if (data.archivos_salida == null) {
                    if (!data.procesado) {
                        mensaje.textContent = data.mensaje
                        mensaje.style.display = 'block'
                        btnEnviar.disabled = false
                        loader.style.display = 'none'
                    }
                } else {
                    data.archivos_salida.forEach(archivo => {
                        var span = document.createElement("span")
                        span.textContent = "Archivo guardado en: " + archivo
                        mensaje.appendChild(span)
                    })
                    // mensaje.textContent = "El archivo se guardó en: " + data["archivo_salida"]
                    mensaje.style.display = 'block'
                    btnEnviar.disabled = false              // Habilito el boton
                    loader.style.display = 'none'           // Saco loader
                }

            })
    })
})