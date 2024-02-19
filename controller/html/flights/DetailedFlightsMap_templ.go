// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package flights

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"github.com/FACorreiaa/Aviation-tracker/controller/models"
)

func detailedMapContainer(data models.LiveFlights) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_detailedMapContainer_a5a3`,
		Function: `function __templ_detailedMapContainer_a5a3(data){console.log('data', data)
    const tileLayer = new ol.layer.Tile({
        // source: new ol.source.StadiaMaps({
        //     layer: 'stamen_toner',
        // }),
        source: new ol.source.OSM(),

    });

    const map = new ol.Map({
        layers: [tileLayer],
        target: 'map',
        view: new ol.View({
            center: ol.proj.fromLonLat([
                parseFloat(data.departure_longitude),
                parseFloat(data.departure_latitude),
            ]),
            zoom: 2,
        }),
    });

    const style = new ol.style.Style({
        stroke: new ol.style.Stroke({
            color: '#EAE911',
            width: 2,
        }),
    });

    const departureMarker = new ol.Feature({
        geometry: new ol.geom.Point(ol.proj.fromLonLat([
            parseFloat(data.departure_longitude),
            parseFloat(data.departure_latitude),
        ])),
        departure: data.departure_airport,
        timezone: data.departure.timezone
    });
    const departureMarkerStyle = new ol.style.Style({
        image: new ol.style.Icon({
            anchor: [0.5, 46],
            anchorXUnits: 'fraction',
            anchorYUnits: 'pixels',
            src: '../../../static/icons/airplane-take-off.svg',
            scale: 0.5,
        }),
    });
    departureMarker.setStyle(departureMarkerStyle);

    const arrivalMarker = new ol.Feature({
        geometry: new ol.geom.Point(ol.proj.fromLonLat([
            parseFloat(data.arrival_longitude),
            parseFloat(data.arrival_latitude),
        ])),
        arrival: data.arrival_airport,
        timezone: data.arrival.timezone
    });
    const arrivalMarkerStyle = new ol.style.Style({
        image: new ol.style.Icon({
            anchor: [0.5, 46],
            anchorXUnits: 'fraction',
            anchorYUnits: 'pixels',
            src: '../../../static/icons/airplane-landing.svg',
            scale: 0.5,
        }),
    });
    arrivalMarker.setStyle(arrivalMarkerStyle);


    const markersSource = new ol.source.Vector({
        features: [departureMarker, arrivalMarker],
    });

    const markersLayer = new ol.layer.Vector({
        source: markersSource,
    });

    map.addLayer(markersLayer);
    map.addLayer(new ol.layer.Vector({
        source: new ol.source.Vector(),
        style: style,
    }));

    const flightsSource = new ol.source.Vector({
        attributions: 'Flight data by ' + '<a href="https://openflights.org/data.html">OpenFlights</a>,',
        loader: function() {
            const arcGenerator = new arc.GreatCircle({
                x: parseFloat(data.departure_longitude),
                y: parseFloat(data.departure_latitude),
            }, {
                x: parseFloat(data.arrival_longitude),
                y: parseFloat(data.arrival_latitude),
            });

            const arcLine = arcGenerator.Arc(100, {
                offset: 10
            });

            const features = [];
            arcLine.geometries.forEach(function(geometry) {
                const line = new ol.geom.LineString(geometry.coords);
                line.transform('EPSG:4326', 'EPSG:3857');

                features.push(
                    new ol.Feature({
                        geometry: line,
                        departure: data.departure_airport,
                        arrival: data.arrival_airport,
                        finished: false,
                    })
                );
            });

            addLater(features, 0);
            tileLayer.on('postrender', animateFlights);
        },
    });

    const flightsLayer = new ol.layer.Vector({
        source: flightsSource,
        style: function(feature) {
            if (feature.get('finished')) {
                return style;
            }
            return null;
        },
    });

    map.addLayer(flightsLayer);

    //popup code
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
            popover.dispose()
            popover = undefined
        } else {
            return
        }
    }

    const tippyButton = document.getElementById('popup');
    tippy(tippyButton, {
        content: document.createElement('div'),
        interactive: true,
        trigger: 'click',
        placement: 'top',
        animation: 'scale',
        theme: 'translucent'
    });

    map.on('click', function(evt) {
        const feature = map.forEachFeatureAtPixel(evt.pixel, function(feature) {
            return feature;
        });
        disposePopover();
        if (!feature) {
            return;
        }
        popup.setPosition(evt.coordinate);

        // Show tooltip with departure or arrival information
        const contentDiv = document.createElement('div');
        if (feature.get('departure')) {
            contentDiv.innerHTML = ` + "`" + `<strong>Departure:</strong> ${feature.get('departure')}<br>
                                <p><strong>Timezone:</strong> ${feature.get('timezone')}<br></p>` + "`" + `;
        } else if (feature.get('arrival')) {
            contentDiv.innerHTML = ` + "`" + `<strong>Arrival:</strong> ${feature.get('arrival')}<br>
                                <p><strong>Timezone:</strong> ${feature.get('timezone')}<br></p>` + "`" + `;
        }
        tippyButton._tippy.setContent(contentDiv);
        tippyButton._tippy.show();
    });


    const pointsPerMs = 0.05;

    function animateFlights(event) {
        const vectorContext = ol.render.getVectorContext(event);
        const frameState = event.frameState;
        vectorContext.setStyle(style);

        const features = flightsSource.getFeatures();
        for (let i = 0; i < features.length; i++) {
            const feature = features[i];
            if (!feature.get('finished')) {
                const coords = feature.getGeometry().getCoordinates();
                const elapsedTime = frameState.time - feature.get('start');
                if (elapsedTime >= 0) {
                    const elapsedPoints = elapsedTime * pointsPerMs;

                    if (elapsedPoints >= coords.length) {
                        feature.set('finished', true);
                    }

                    const maxIndex = Math.min(elapsedPoints, coords.length);
                    const currentLine = new ol.geom.LineString(coords.slice(0, maxIndex));

                    const worldWidth = ol.extent.getWidth(map.getView().getProjection().getExtent());
                    const offset = Math.floor(map.getView().getCenter()[0] / worldWidth);

                    currentLine.translate(offset * worldWidth, 0);
                    vectorContext.drawGeometry(currentLine);
                    currentLine.translate(worldWidth, 0);
                    vectorContext.drawGeometry(currentLine);
                }
            }
        }
        map.render();
    }

    function addLater(features, timeout) {
        window.setTimeout(function() {
            let start = Date.now();
            features.forEach(function(feature) {
                feature.set('start', start);
                flightsSource.addFeature(feature);
                const duration = (feature.getGeometry().getCoordinates().length - 1) / pointsPerMs;
                start += duration;
            });
        }, timeout);
    }
}`,
		Call:       templ.SafeScript(`__templ_detailedMapContainer_a5a3`, data),
		CallInline: templ.SafeScriptInline(`__templ_detailedMapContainer_a5a3`, data),
	}
}

func FlightsDetailMap(data models.LiveFlights) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<html><head><style scoped>\n                .map {\n                    height: 200px;\n                    z-index:1;\n                }\n\n                .map-container {\n                    position: relative;\n                    min-height: 400px;\n                    padding-top: 10px;\n                    z-index: 0;\n                }\n\n                a.skiplink {\n                    position:\n                        absolute clip: rect(1px, 1px, 1px, 1px);\n                    padding:\n                        0 border: 0 height: 1px;\n                    width: 1px;\n                    overflow:\n                        hidden\n                }\n\n                a.skiplink:focus {\n                    clip:\n                        auto height: auto width: auto background-color: #fff;\n                    padding: 0.3em;\n                }\n\n                #map:focus {\n                    outline: #4A74A8 solid 0.15em;\n                }\n\t\t\t</style></head>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.RenderScriptItems(ctx, templ_7745c5c3_Buffer, detailedMapContainer(data))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<body onload=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 templ.ComponentScript = detailedMapContainer(data)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var2.Call)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><div class=\"container map-container\"><div id=\"map\" class=\"mx-auto map\" tabindex=\"0\"><button aria-describedby=\"popup\" data-tippy-content=\"popup\" id=\"popup\"></button></div></div><button id=\"zoom-out\" class=\"btn btn-secondary\">Zoom out</button> <button id=\"zoom-in\" class=\"btn btn-secondary\">Zoom in</button><script src=\"https://api.mapbox.com/mapbox.js/plugins/arc.js/v0.1.0/arc.js\"></script></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}