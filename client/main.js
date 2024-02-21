console.log("prueba")
document.addEventListener("DOMContentLoaded", function () {
  const select = document.getElementById("select");

  // Realizar la solicitud de red (fetch) al servidor
  fetch("/procesos")
      .then(response => {
        // response.json().then(console.log(...data))
          if (!response.ok) {
              throw new Error('La solicitud falló: ' + response.status);
          }
          return response.json();
      })
      .then(data => {
          // Limpiar opciones existentes
          select.innerHTML = "";

          // Agregar nuevas opciones al select
          data.forEach(proceso => {
              const option = document.createElement("option");
              option.value = proceso.id; // Asegúrate de que proceso tenga una propiedad "id"
              option.textContent = proceso.nombre; // Asegúrate de que proceso tenga una propiedad "nombre"
              select.appendChild(option);
          });
      })
      .catch(error => {
          console.error("Error al obtener los procesos:", error);
      });
});