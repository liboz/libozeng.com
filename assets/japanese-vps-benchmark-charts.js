luxon.Settings.defaultLocale = "ja-JP";

let mapper = {
  filterFunction: (r) => {
    return (
      r["ping_test_results"].filter((i) => i.url === "google.com").length > 0
    );
  },
  getterFunction: (r) => {
    return {
      x: r["start_time"] * 1000,
      y: r["ping_test_results"].filter((i) => i.url === "google.com")[0][
        "average_ping"
      ],
    };
  },
  generateDataset: () => {
    return Object.entries(sqlData).map((kvp) => {
      const name = kvp[0];
      const value = kvp[1];
      return {
        backgroundColor: colorMap[name],
        label: name,
        data: value.filter(mapper.filterFunction).map(mapper.getterFunction),
      };
    });
  },
  labelFunction: (context) => {
    return (
      context.dataset.label +
      " - " +
      context.label +
      " - " +
      context.raw.y +
      "ms"
    );
  },
  yScaleFunction: (value, index, ticks) => {
    return value + "ms";
  },
};

function changeChart() {
  const primaryValue = document.getElementById("primarySelect").value;
  if (primaryValue.includes("ping_test_results")) {
    document.getElementById("siteSelect").style = "";
  } else {
    document.getElementById("siteSelect").style = "display: none;";
  }

  if (primaryValue.includes("speed_test_results")) {
    document.getElementById("locationSelect").style = "";
  } else {
    document.getElementById("locationSelect").style = "display: none;";
  }

  if (primaryValue.includes("ping_test_results")) {
    const [firstKey, secondKey] = primaryValue.split(",");
    const site = document.getElementById("siteSelect").value;

    mapper.getterFunction = (r) => {
      return {
        x: r["start_time"] * 1000,
        y: r[firstKey].filter((i) => i.url === site)[0][secondKey],
      };
    };

    mapper.filterFunction = (r) => {
      return r[firstKey].filter((i) => i.url === site).length > 0;
    };
    mapper.labelFunction = (context) => {
      return (
        context.dataset.label +
        " - " +
        context.label +
        " - " +
        context.raw.y +
        "ms"
      );
    };
    mapper.yScaleFunction = (value, index, ticks) => {
      return value + "ms";
    };
  } else if (primaryValue.includes("speed_test_results")) {
    const [firstKey, secondKey] = primaryValue.split(",");
    const location = document.getElementById("locationSelect").value;
    mapper.getterFunction = (r) => {
      return {
        x: r["start_time"] * 1000,
        y: r[firstKey]?.filter((i) => i.city === location)[0][secondKey],
      };
    };
    mapper.filterFunction = (r) => {
      return r[firstKey]?.filter((i) => i.city === location).length > 0;
    };

    mapper.labelFunction = (context) => {
      return (
        context.dataset.label +
        " - " +
        context.label +
        " - " +
        formatBytes(context.raw.y) +
        "/s"
      );
    };
    mapper.yScaleFunction = (value, index, ticks) => {
      return formatBytes(value) + "/s";
    };
  } else {
    mapper.getterFunction = (r) => {
      return {
        x: r["start_time"] * 1000,
        y: r[primaryValue],
      };
    };
    mapper.filterFunction = (r) => true;

    if (primaryValue.includes("geekbench")) {
      mapper.labelFunction = (context) => {
        return (
          context.dataset.label + " - " + context.label + " - " + context.raw.y
        );
      };
      mapper.yScaleFunction = (value, index, ticks) => {
        return value;
      };
    } else {
      mapper.labelFunction = (context) => {
        return (
          context.dataset.label +
          " - " +
          context.label +
          " - " +
          formatBytes(context.raw.y) +
          "/s"
        );
      };
      mapper.yScaleFunction = (value, index, ticks) => {
        return formatBytes(value) + "/s";
      };
    }
  }
  myChart.data.datasets = mapper.generateDataset();
  myChart.options.plugins.tooltip.callbacks.label = mapper.labelFunction;
  myChart.options.scales.y.ticks.callback = mapper.yScaleFunction;
  myChart.update();
  console.log("test", primaryValue);
}

function formatBytes(bytes, decimals = 2) {
  if (!+bytes) return "0 Bytes";

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
}

const colorMap = {
  xserver: "#CFFFB3",
  vultr: "#ADE25D",
  linode: "#FCEC52",
  kagoya: "#0AD3FF",
  conoha: "#3B7080",
  amazon: "#3A5743",
};

const ctx = document.getElementById("myChart").getContext("2d");
const myChart = new Chart(ctx, {
  type: "scatter",
  data: {
    datasets: Object.entries(sqlData).map((kvp) => {
      const name = kvp[0];
      const value = kvp[1];
      return {
        backgroundColor: colorMap[name],
        label: name,
        data: value.filter(mapper.filterFunction).map(mapper.getterFunction),
      };
    }),
  },
  options: {
    plugins: {
      tooltip: {
        callbacks: {
          label: mapper.labelFunction,
        },
      },
    },
    parsing: false,
    scales: {
      x: {
        type: "time",
        time: {
          unit: "day",
        },
      },
      y: {
        type: "logarithmic",
        ticks: {
          callback: mapper.yScaleFunction,
        },
      },
    },
  },
});
