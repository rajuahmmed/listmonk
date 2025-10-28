<template>
  <div class="items">
    <b-tabs :animated="false" v-model="tab">
      <b-tab-item :label="$t('settings.appearance.adminName')" label-position="on-border">
        <div class="block">
          {{ $t('settings.appearance.adminHelp') }}
        </div>

        <b-field :label="$t('settings.appearance.customCSS')" label-position="on-border">
          <code-editor lang="css" v-model="data['appearance.admin.custom_css']" name="body" key="editor-admin-css" />
        </b-field>

        <b-field :label="$t('settings.appearance.customJS')" label-position="on-border">
          <code-editor lang="javascript" v-model="data['appearance.admin.custom_js']" name="body"
            key="editor-admin-js" />
        </b-field>
      </b-tab-item><!-- admin -->

      <b-tab-item :label="$t('settings.appearance.publicName')" label-position="on-border">
        <div class="block">
          {{ $t('settings.appearance.publicHelp') }}
        </div>

        <b-field :label="$t('settings.appearance.customCSS')" label-position="on-border">
          <code-editor lang="css" v-model="data['appearance.public.custom_css']" name="body" key="editor-public-css" />
        </b-field>

        <b-field :label="$t('settings.appearance.customJS')" label-position="on-border">
          <code-editor lang="javascript" v-model="data['appearance.public.custom_js']" name="body"
            key="editor-public-js" />
        </b-field>
      </b-tab-item><!-- public -->

      <b-tab-item :label="$t('settings.appearance.buttonStyleName')" label-position="on-border">
        <div class="block">
          {{ $t('settings.appearance.buttonStyleHelp') }}
        </div>

        <div class="columns">
          <div class="column is-6">
            <b-field :label="$t('settings.appearance.buttonBgColor')" label-position="on-border">
              <b-input v-model="data['appearance.button.bg_color']"
                       type="text"
                       placeholder="#0055d4" />
            </b-field>
          </div>
          <div class="column is-6">
            <b-field :label="$t('settings.appearance.buttonTextColor')" label-position="on-border">
              <b-input v-model="data['appearance.button.text_color']"
                       type="text"
                       placeholder="#ffffff" />
            </b-field>
          </div>
        </div>

        <div class="columns">
          <div class="column is-6">
            <b-field :label="$t('settings.appearance.buttonHoverBgColor')" label-position="on-border">
              <b-input v-model="data['appearance.button.hover_bg_color']"
                       type="text"
                       placeholder="#222222" />
            </b-field>
          </div>
          <div class="column is-6">
            <b-field :label="$t('settings.appearance.buttonHoverTextColor')" label-position="on-border">
              <b-input v-model="data['appearance.button.hover_text_color']"
                       type="text"
                       placeholder="#ffffff" />
            </b-field>
          </div>
        </div>

        <div class="columns">
          <div class="column is-6">
            <b-field :label="$t('settings.appearance.buttonBorderRadius')" label-position="on-border">
              <b-input v-model="data['appearance.button.border_radius']"
                       type="text"
                       placeholder="3px" />
            </b-field>
          </div>
        </div>

        <div class="block">
          <p class="is-size-7 has-text-grey">
            {{ $t('settings.appearance.buttonStylePreview') }}
          </p>
          <a :style="{
            background: data['appearance.button.bg_color'] || '#0055d4',
            color: data['appearance.button.text_color'] || '#ffffff',
            borderRadius: data['appearance.button.border_radius'] || '3px',
            padding: '10px 30px',
            display: 'inline-block',
            textDecoration: 'none',
            fontWeight: 'bold',
            marginTop: '10px',
          }" href="#" @click.prevent>
            {{ $t('settings.appearance.buttonStylePreviewText') }}
          </a>
        </div>
      </b-tab-item><!-- button style -->
    </b-tabs>
  </div>
</template>

<script>
import Vue from 'vue';
import { mapState } from 'vuex';
import CodeEditor from '../../components/CodeEditor.vue';

export default Vue.extend({
  components: {
    'code-editor': CodeEditor,
  },

  props: {
    form: {
      type: Object, default: () => { },
    },
  },

  data() {
    return {
      data: this.form,
      tab: 0,
    };
  },

  mounted() {
    this.tab = this.$utils.getPref('settings.apperanceTab') || 0;
  },

  watch: {
    tab(t) {
      this.$utils.setPref('settings.apperanceTab', t);
    },
  },

  computed: {
    ...mapState(['settings']),
  },
});

</script>
