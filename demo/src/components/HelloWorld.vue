<template>
  <h1>🚀 {{ msg }} 🌏</h1>

  <div v-if="danger != ''" class="w3-panel w3-pale-red w3-border">
    <h3>{{ danger }}</h3>
  </div>

  <div v-if="warning != ''" class="w3-panel w3-pale-yellow w3-border">
    <h3>{{ warning }}</h3>
  </div>

  <div
    class="w3-container w3-margin"
    :class="{ 'w3-disabled': danger != '' || finish }"
  ></div>
  <template> </template>
  <!-- 渲染每個步驟 -->
  <div
    v-for="(stage, index_stage) in show_stages"
    :key="index_stage"
    class="w3-center w3-container w3-margin"
    :class="{
      loading,
      'w3-disabled': finish,
      'w3-panel w3-leftbar w3-border-blue': step == index_stage + 1,
    }"
  >
    <!-- 步驟內容 -->
    <h3>{{ index_stage + 1 }}. {{ stage.title }}</h3>
    <template v-if="stage.no_action">
      <span v-if="stage.done"> ✅ 完成</span>
      <span v-else>⏱ 執行中</span>
    </template>
    <!-- 多選 NODE -->
    <template v-else-if="stage.more">
      <template v-for="(node, index_node) in nodes" :key="index_node">
        <div
          v-if="node_master === node"
          class="disabled w3-margin-left w3-btn w3-border w3-border-blue w3-blue"
        >
          {{ node }}
        </div>
        <div
          v-else-if="!in_server_node(index_stage, node)"
          class="w3-margin-left w3-btn w3-border w3-border-blue"
          @click="toggle_in_stage(index_stage, node)"
        >
          {{ node }}
        </div>
        <div
          v-else
          class="w3-margin-left w3-btn w3-border w3-border-blue w3-blue"
          :class="{
            loading,
          }"
          @click="toggle_in_stage(index_stage, node)"
        >
          {{ node }}
        </div>
      </template>
      <div
        class="w3-margin-left w3-btn w3-border w3-border-blue"
        @click="stage.decide(index_stage)"
      >
        確認
      </div>
    </template>
    <!-- 單選 NODE -->
    <template v-else>
      <template v-for="(node, index_node) in nodes" :key="index_node">
        <div
          v-if="!has_decide(index_stage)"
          class="w3-margin-left w3-btn w3-border w3-border-blue"
          @click="stage.decide(index_stage, node)"
        >
          {{ node }}
        </div>
        <div
          v-else
          class="w3-margin-left w3-btn w3-border w3-border-blue"
          :class="{
            'w3-blue': node == stage.selected_node,
            loading,
            disabled: !loading,
          }"
        >
          {{ node }}
        </div>
      </template>
    </template>
  </div>
  <!-- 完成步驟 -->
  <h1 v-if="finish">🥳 Clsuter Setup Finish 🍺</h1>
</template>

<script>
export default {
  name: "HelloWorld",
  props: {
    msg: String,
  },
  data() {
    return {
      ws: null,
      danger: "",
      warning: "",
      stages: [],
      nodes: [],
      uuid: "",
      node_master: "",
      step: 0,
      loading: false,
      finish: false,
    };
  },
  computed: {
    show_stages() {
      if (this.step < 1) {
        return [];
      }
      return this.stages.slice(0, this.step);
    },
  },
  methods: {
    reset() {
      this.warning = "";
      this.nodes = [];
      this.uuid = "";
      this.step = 0;
      this.loading = false;
    },
    send(data) {
      if (this.ws) {
        this.ws.send(
          JSON.stringify({
            node: data,
            uuid: this.uuid,
          })
        );
      } else {
        this.danger = "請重整頁面，重新連線";
      }
    },
    toggle_in_stage(index_stage, node) {
      let stage = this.stages[index_stage];
      if (stage.more) {
        if (node == this.node_master) {
          return;
        }

        const index = stage.selected_node.indexOf(node);
        if (index > -1) {
          stage.selected_node = stage.selected_node.filter((n) => {
            return node != n;
          });
        } else {
          stage.selected_node.push(node);
        }

        if (stage.selected_node.length % 2 == 0) {
          this.warning = "";
        }
      } else {
        stage.selected_node = stage.selected_node === node ? "" : node;
      }
    },
    in_server_node(index_stage, node) {
      let stage = this.stages[index_stage];
      if (!stage.more) {
        return false;
      }
      return stage.selected_node.indexOf(node) > -1;
    },
    has_decide(index_stage) {
      let stage = this.stages[index_stage];
      return this.nodes.indexOf(stage.selected_node) > -1;
    },
    decide_node(index_stage, n) {
      this.decide(n);
      this.stages[index_stage].selected_node = n;
    },
    decide_more(index_stage) {
      if (this.stages[index_stage].selected_node.length % 2 == 1) {
        this.warning = `Server Node 需為奇數 (目前 ${
          this.stages[index_stage].selected_node.length + 1
        })`;
        return;
      }
      this.decide(this.stages[index_stage].selected_node);
    },
    decide(n) {
      this.loading = true;
      this.send(n);
    },
    open_conn(ws) {
      this.ws = ws;
      this.danger = "";
    },
    conn_closed() {
      this.reset();
    },
    process_error(e) {
      this.danger = e;
    },
    receive_message(e) {
      let data = JSON.parse(e.data);
      // console.log("接收資料", data);

      if (data.Event === "get_ip") {
        this.nodes = data.Nodes;
      }

      if (data.Step === -1) {
        return;
      }

      this.loading = false;
      this.uuid = data.UUID;
      this.step = data.Step;
      const index_stage = data.Step - 1;
      setTimeout(() => window.scrollBy(0, document.body.scrollHeight), 200);

      if (data.Event === "finish") {
        this.finish = true;
        setTimeout(() => window.scrollBy(0, document.body.scrollHeight), 200);
        return;
      }

      if (this.stages[index_stage].no_action) {
        if (data.Event === "done") {
          this.stages[index_stage].done = true;
          setTimeout(() => window.scrollBy(0, document.body.scrollHeight), 200);
        } else if (data.Event === "start") {
          this.stages[index_stage].done = false;
          setTimeout(() => window.scrollBy(0, document.body.scrollHeight), 200);
        }
        return;
      }

      if (data.Event === "get_ip") {
        if (this.stages[index_stage].more) {
          this.stages[index_stage].selected_node = [];
        } else {
          this.stages[index_stage].selected_node = "";
        }
      }
    },
  },
  mounted() {
    // console.log("mounted");
    this.stages = [
      {
        title: "請選擇一台 Registry Node",
        selected_node: "",
        decide: this.decide_node,
      },
      {
        title: "請選擇一台 Master Node",
        selected_node: "",
        decide: (index_stage, n) => {
          this.decide_node(index_stage, n);
          this.node_master = n;
        },
      },
      {
        title: "請選擇 Server Node (未選作為 Agent Node)",
        selected_node: [],
        more: true,
        decide: this.decide_more,
      },
      {
        title: "請選擇一台 Node 作為日誌儲存",
        selected_node: "",
        decide: this.decide_node,
      },
      {
        title: "請選擇一台 Node 作為監控儲存",
        selected_node: "",
        decide: this.decide_node,
      },
      {
        title: "請選擇一台 Node 作為監控介面(Dashboard)",
        selected_node: "",
        decide: this.decide_node,
      },
      {
        title: "部署Database",
        done: false,
        no_action: true,
      },
      {
        title: "部署AI",
        done: false,
        no_action: true,
      },
      {
        title: "部署排程",
        done: false,
        no_action: true,
      },
      {
        title: "部署服務",
        done: false,
        no_action: true,
      },
    ];
    connect(
      this.open_conn,
      this.receive_message,
      this.process_error,
      this.conn_closed
    );
  },
};

let connect = (open_conn, receive_message, process_error, conn_closed) => {
  let protocol = "ws";
  if (window.location.protocol === "https:") {
    protocol = "wss";
  }
  const url = `${protocol}://${window.location.host}/deploy/guide`;
  // console.log(url);
  let ws = new WebSocket(url);

  ws.onopen = () => {
    console.info("開始連線");
    open_conn(ws);
  };
  ws.onmessage = (e) => {
    // console.log("接收訊息", e);
    receive_message(e);
  };
  ws.onerror = (e) => {
    console.error("連線錯誤", e);
    process_error("連線失敗，請刷新頁面");
  };
  ws.onclose = (e) => {
    console.warn("連線關閉", e);
    conn_closed();
    setTimeout(
      () => connect(open_conn, receive_message, process_error, conn_closed),
      2000
    );
  };
};
</script>

<style>
.loading {
  cursor: wait;
}

.disabled {
  cursor: not-allowed;
}
</style>
