<template>
<div>
    <div class="card contain-card text-black-50 bg-white mb-5">
        <div class="card-header">{{message.id ? `Update ${message.title}` : "Create Annoucement"}}

            <transition name="slide-fade">
                <button @click="removeEdit" v-if="message.id" class="btn btn-sm float-right btn-danger btn-sm">Close</button>
            </transition>
        </div>
        <div class="card-body">
    <form @submit="saveMessage">
        <div class="form-group row">
            <label class="col-sm-4 col-form-label">Title</label>
            <div class="col-sm-8">
                <input v-model="message.title" type="text" name="title" class="form-control" id="title" placeholder="Announcement Title" required>
            </div>
        </div>

        <div class="form-group row">
            <label class="col-sm-4 col-form-label">Description</label>
            <div class="col-sm-8">
                <textarea v-model="message.description" rows="5" name="description" class="form-control" id="description" required></textarea>
            </div>
        </div>

        <div class="form-group row">
            <label class="col-sm-4 col-form-label">Service</label>
            <div class="col-sm-8">
                <select v-model="message.service_id" name="service_id" class="form-control">
                    <option :value="0">Global Announcement</option>
                    <option v-for="(service, i) in $store.getters.services" :value="service.id" >{{service.name}}</option>
                </select>
            </div>
        </div>

        <div class="form-group row">
            <label class="col-sm-4 col-form-label">Announcement Date Range</label>
            <div class="col-sm-4">
                <flatPickr v-model="message.start_on" @on-change="startChange" :config="config" type="text" name="start_on" class="form-control form-control-plaintext" id="start_on" value="0001-01-01T00:00:00Z" required />
            </div>
            <div class="col-sm-4">
                <flatPickr v-model="message.end_on" @on-change="endChange" :config="config" type="text" name="end_on" class="form-control form-control-plaintext" id="end_on" value="0001-01-01T00:00:00Z" required />
            </div>
        </div>

        <div v-show="this.service === null" class="form-group row">
            <label for="service_id" class="col-sm-4 col-form-label">Service</label>
            <div class="col-sm-8">
                <select v-model="message.service_id" class="form-control" name="service" id="service_id">
                    <option :value="0">Global Message</option>
                    <option v-for="(service, index) in $store.getters.services" :value="service.id" v-bind:key="index" >{{service.name}}</option>
                </select>
            </div>
        </div>

        <div class="form-group row">
            <label for="notify_method" class="col-sm-4 col-form-label">Notify Users</label>
            <div class="col-sm-8">
                <span @click="message.notify = !!message.notify" class="switch">
                    <input v-model="message.notify" type="checkbox" class="switch" id="switch-normal">
                    <label for="switch-normal">Notify Users Before Scheduled Time</label>
                </span>
            </div>
        </div>

        <div v-if="message.notify" class="form-group row">
            <label for="notify_method" class="col-sm-4 col-form-label">Notification Method</label>
            <div class="col-sm-8">
                <input v-model="message.notify_method" type="text" name="notify_method" class="form-control" id="notify_method" value="" placeholder="email">
            </div>
        </div>

        <div v-if="message.notify" class="form-group row">
            <label for="notify_before" class="col-sm-4 col-form-label">Notify Before</label>
            <div class="col-sm-8">
                <div class="form-inline">
                    <input v-model="message.notify_before" type="number" name="notify_before" class="col-4 form-control" id="notify_before">
                    <select v-model="message.notify_before_scale" class="ml-2 col-7 form-control" name="notify_before_scale" id="notify_before_scale">
                        <option value="minute">Minutes</option>
                        <option value="hour">Hours</option>
                        <option value="day">Days</option>
                    </select>
                </div>
            </div>
        </div>

        <div class="form-group row">
            <div class="col-sm-12">
                <button @click="saveMessage"
                        :disabled="!message.title || !message.description"
                        type="submit" class="btn btn-block" :class="{'btn-primary': !message.id, 'btn-secondary': message.id}">
                    {{message.id ? "Edit Message" : "Create Message"}}
                </button>
            </div>
        </div>
        <div class="alert alert-danger d-none" id="alerter" role="alert"></div>
    </form>
            </div>
    </div>
</div>
</template>

<script>
  import Api from "../API";
  import flatPickr from 'vue-flatpickr-component';
  import 'flatpickr/dist/flatpickr.css';

  export default {
  name: 'FormMessage',
  components: {
    flatPickr
  },
  props: {
    in_message: {
      type: Object
    },
      service: {
          type: Object
      },
    edit: {
      type: Function
    }
  },
  data () {
    return {
      message: {
        title: "",
        description: "",
        start_on: new Date(),
        end_on: new Date(),
        service_id: 0,
        notify_method: "",
        notify: false,
        notify_before: 0,
        notify_before_scale: "minute",
      },
      config: {
        altFormat: "l M J, \\at h:iK",
        altInput: true,
        enableTime: true,
        dateFormat: "Z",
      },
      temp: {}
    }
  },
  watch: {
    in_message() {
      this.message = this.in_message
    }
  },
      mounted () {
        if (this.service) {
            this.service_id = this.service.id
        }
      },
      methods: {
      startChange(e) {
        window.console.log(e)
      },
      endChange(e) {
          window.console.log(e)
      },
    removeEdit() {
      this.message = {}
      this.edit(false)
    },
    async saveMessage(e) {
      e.preventDefault();
      if (this.message.id) {
        await this.updateMessage()
      } else {
        await this.createMessage()
      }
    },
    async createMessage() {
      await Api.message_create(this.message)
      const messages = await Api.messages()
      this.$store.commit('setMessages', messages)
      this.message = {}
    },
    async updateMessage() {
      await Api.message_update(this.message)
      const messages = await Api.messages()
      this.$store.commit('setMessages', messages)
      this.edit(false)
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
