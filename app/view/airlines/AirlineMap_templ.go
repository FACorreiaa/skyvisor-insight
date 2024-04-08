// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.648
package airline

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"fmt"
	"github.com/FACorreiaa/Aviation-tracker/app/models"
)

func mapContainer(data []models.Airline) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_mapContainer_87ed`,
		Function: `function __templ_mapContainer_87ed(data){//control selector
            const rangeInput = document.querySelector('.range');
            const updatePointsOnMap = () => {
                const selectedValue = parseInt(rangeInput.value, 10);
                // Logic to update the number of points on the map based on the selected value
                const filteredData = data.slice(0, selectedValue);
                // Clear existing features
                vectorSource.clear();
                // Add new features based on the filtered data
                vectorSource.addFeatures(filteredData.map(airline => createFeatureFromAirline(airline)));
            };

        // Add event listener for input change
        rangeInput.addEventListener('input', updatePointsOnMap);

  function createFeatureFromAirline(airline) {
        const iconStyle = new ol.style.Style({
            image: new ol.style.Icon({
                anchor: [1, 46],
                anchorXUnits: 'fraction',
                anchorYUnits: 'pixels',
                src: '../../static/icons/marker.png',
                scale: 0.5,
            }),
        });

        const feature = new ol.Feature({
            geometry: new ol.geom.Point(ol.proj.fromLonLat([airline.longitude, airline.latitude])),
            airline: airline.airline_name,
            airport: airline.airport_name,
            fleet_average: airline.fleet_average_age,
            fleet_size: airline.fleet_size,
            country: airline.country_name,
            city: airline.city_name,
            status: airline.status,
            style: iconStyle,
        });

        // Set the minimum and maximum zoom levels for the marker to be visible
        feature.set('minZoom', 10); // Adjust the zoom level as needed
        feature.set('maxZoom', 18); // Adjust the zoom level as needed

        feature.setStyle(iconStyle);

        return feature;
    }


    const vectorSource = new ol.source.Vector({
        features: data.map(airline => createFeatureFromAirline(airline)),
    });

   const vectorLayer = new ol.layer.Vector({
      source: vectorSource,
   });

  const tileLayer = new ol.layer.Tile({
            source: new ol.source.OSM(),
         })

   const map = new ol.Map({
      layers: [tileLayer, vectorLayer],
      target: document.getElementById('map'),
      controls: [],
      view: new ol.View({
         center: [0, 0],
         zoom: 1,
      }),
   });

   const element = document.getElementById('popup');
   const popup = new ol.Overlay({
      element: element,
      positioning: 'bottom-center',
      stopEvent: false,
   });
   map.addOverlay(popup);

   let popover;

   function disposePopover() {
      if (popover) {
         popover.dispose();
         popover = undefined;
      }
   }

   const tippyButton = document.getElementById('popup');
    tippy(tippyButton, {
      content: document.createElement('div'),
      interactive: true,
      trigger: 'click',
      placement: 'top',
      animation: 'scale'  ,
      theme: 'translucent'
    });

   map.on('click', function (evt) {
      const feature = map.forEachFeatureAtPixel(evt.pixel, function (feature) {
         return feature;
      });
      disposePopover();
      if (!feature) {
         return;
      }
      popup.setPosition(evt.coordinate);

// airline: airport.airline_name
//             airport: airport.airport_name,
//             fleet_average: airport.fleet_average_age,
//             fleet_size: airport.fleet_size,
//             country: airport.country_name,
//             city: airport.city_name,
//             timezone: airport.timezone,
//             status: airport.status,

         //show toolip
      const contentDiv = document.createElement('div');
      contentDiv.innerHTML = ` + "`" + `
            <strong>Airline:</strong> ${feature.get('airline')}<br>
            <strong>Airport:</strong> ${feature.get('airport')}<br>
            <strong>Fleet average size:</strong> ${feature.get('fleet_average')}<br>
            <strong>Fleet size:</strong> ${feature.get('fleet_size')}<br>
            <strong>Country:</strong> ${feature.get('country')}<br>
            <strong>City:</strong> ${feature.get('city')}<br>
            <strong>Status:</strong> ${feature.get('status')}<br>

    ` + "`" + `;
      tippyButton._tippy.setContent(contentDiv);
      tippyButton._tippy.show();
   });



   map.on('pointermove', function (e) {
      const pixel = map.getEventPixel(e.originalEvent);
      const hit = map.hasFeatureAtPixel(pixel);
      map.getTarget().style.cursor = hit ? "pointer" : "";
   });


   map.on('movestart', disposePopover);


   document.getElementById('zoom-out').onclick = function () {
      const view = map.getView();
      const zoom = view.getZoom();
      view.setZoom(zoom - 1);
   };

   document.getElementById('zoom-in').onclick = function () {
      const view = map.getView();
      const zoom = view.getZoom();
      view.setZoom(zoom + 1);
   };

   map.on('dblclick', event => {
       const zoomLevel = 8;
       const feature = map.forEachFeatureAtPixel(event.pixel, (feature) => {
           return feature;
       });
       if (feature instanceof ol.Feature) {
           map.getView().fit(feature.getGeometry().getExtent(), {
               size: map.getSize(),
               padding: [10, 10, 10, 10],
               minResolution: map.getView().getResolutionForZoom(zoomLevel),
           });
       }
   });

   map.getView().on('change:resolution', function () {
               if (map.getView().getZoom() < 4) {
                   vectorLayer.setVisible(false);
               } else {
                   vectorLayer.setVisible(true);
               }
           });

   map.on('postrender', function () {
                     if (map.getView().getZoom() < 3) {
                         vectorLayer.setVisible(false);
                     } else {
                         vectorLayer.setVisible(true);
                     }
                 });
}`,
		Call:       templ.SafeScript(`__templ_mapContainer_87ed`, data),
		CallInline: templ.SafeScriptInline(`__templ_mapContainer_87ed`, data),
	}
}

func AirlineMap(data []models.Airline) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<html><head><style scoped>\n\t\t\t\t.map {\n\t\t\t\t\theight: 600px;\n                    z-index:1;\n                }\n\t\t\t\t#map:focus {\n\t\t\t\t\toutline: #4A74A8 solid 0.15em;\n\t\t\t\t}\n\t\t\t</style></head>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.RenderScriptItems(ctx, templ_7745c5c3_Buffer, mapContainer(data))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<body onload=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 templ.ComponentScript = mapContainer(data)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var2.Call)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><div><div class=\"mb-5 text-left\"><button id=\"zoom-out\" class=\"btn btn-ghost mr-5\">Zoom out</button> <button id=\"zoom-in\" class=\"btn btn-ghost\">Zoom in</button></div><div id=\"map\" class=\" map\" tabindex=\"0\"><button aria-describedby=\"popup\" data-tippy-content=\"popup\" id=\"popup\"></button></div><div class=\"w-full form-control\"><div class=\"label\"><span class=\"label-text font-semi-bold text-xs badge-xs badge pb-0\">Display number of markers</span></div><div class=\"flex items-center mt-2\"><span class=\"text-xs\">0</span> <input type=\"range\" min=\"0\" max=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%d", len(data)))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/airlines/AirlineMap.templ`, Line: 228, Col: 70}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%d", len(data)))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/airlines/AirlineMap.templ`, Line: 229, Col: 72}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" class=\"range range-xs p-2 mr-2 ml-2\" name=\"rangeValue\"> <span class=\"text-xs\" id=\"rangeValue\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var5 string
		templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%d", len(data)))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/airlines/AirlineMap.templ`, Line: 233, Col: 103}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></div></div></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
