import { useEffect, useMemo, useState } from "react";
import { getI18n, useTranslation } from "react-i18next";
import { type FlowNodeEntity, getNodeForm } from "@flowgram.ai/fixed-layout-editor";
import { IconPlus } from "@tabler/icons-react";
import { type AnchorProps, Button, Divider, Flex, Form, type FormInstance, Select, Switch, Typography, theme } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import AccessEditDrawer from "@/components/access/AccessEditDrawer";
import AccessSelect from "@/components/access/AccessSelect";
import DeploymentProviderPicker from "@/components/provider/DeploymentProviderPicker";
import DeploymentProviderSelect from "@/components/provider/DeploymentProviderSelect";
import Show from "@/components/Show";
import { ACCESS_USAGES, DEPLOYMENT_PROVIDERS, accessProvidersMap, deploymentProvidersMap } from "@/domain/provider";
import { type WorkflowNodeConfigForDeploy, WorkflowNodeType, defaultNodeConfigForDeploy } from "@/domain/workflow";
import { useAntdForm, useZustandShallowSelector } from "@/hooks";
import { useWorkflowStore } from "@/stores/workflow";

import { FormNestedFieldsContextProvider, NodeFormContextProvider } from "./_context";
import BizDeployNodeConfigFormProvider1PanelConsole from "./BizDeployNodeConfigFormProvider1PanelConsole";
import BizDeployNodeConfigFormProvider1PanelSite from "./BizDeployNodeConfigFormProvider1PanelSite";
import BizDeployNodeConfigFormProviderAliyunALB from "./BizDeployNodeConfigFormProviderAliyunALB";
import BizDeployNodeConfigFormProviderAliyunAPIGW from "./BizDeployNodeConfigFormProviderAliyunAPIGW";
import BizDeployNodeConfigFormProviderAliyunCAS from "./BizDeployNodeConfigFormProviderAliyunCAS";
import BizDeployNodeConfigFormProviderAliyunCASDeploy from "./BizDeployNodeConfigFormProviderAliyunCASDeploy";
import BizDeployNodeConfigFormProviderAliyunCDN from "./BizDeployNodeConfigFormProviderAliyunCDN";
import BizDeployNodeConfigFormProviderAliyunCLB from "./BizDeployNodeConfigFormProviderAliyunCLB";
import BizDeployNodeConfigFormProviderAliyunDCDN from "./BizDeployNodeConfigFormProviderAliyunDCDN";
import BizDeployNodeConfigFormProviderAliyunDDoS from "./BizDeployNodeConfigFormProviderAliyunDDoS";
import BizDeployNodeConfigFormProviderAliyunESA from "./BizDeployNodeConfigFormProviderAliyunESA";
import BizDeployNodeConfigFormProviderAliyunFC from "./BizDeployNodeConfigFormProviderAliyunFC";
import BizDeployNodeConfigFormProviderAliyunGA from "./BizDeployNodeConfigFormProviderAliyunGA";
import BizDeployNodeConfigFormProviderAliyunLive from "./BizDeployNodeConfigFormProviderAliyunLive";
import BizDeployNodeConfigFormProviderAliyunNLB from "./BizDeployNodeConfigFormProviderAliyunNLB";
import BizDeployNodeConfigFormProviderAliyunOSS from "./BizDeployNodeConfigFormProviderAliyunOSS";
import BizDeployNodeConfigFormProviderAliyunVOD from "./BizDeployNodeConfigFormProviderAliyunVOD";
import BizDeployNodeConfigFormProviderAliyunWAF from "./BizDeployNodeConfigFormProviderAliyunWAF";
import BizDeployNodeConfigFormProviderAPISIX from "./BizDeployNodeConfigFormProviderAPISIX";
import BizDeployNodeConfigFormProviderAWSACM from "./BizDeployNodeConfigFormProviderAWSACM";
import BizDeployNodeConfigFormProviderAWSCloudFront from "./BizDeployNodeConfigFormProviderAWSCloudFront";
import BizDeployNodeConfigFormProviderAWSIAM from "./BizDeployNodeConfigFormProviderAWSIAM";
import BizDeployNodeConfigFormProviderAzureKeyVault from "./BizDeployNodeConfigFormProviderAzureKeyVault";
import BizDeployNodeConfigFormProviderBaiduCloudAppBLB from "./BizDeployNodeConfigFormProviderBaiduCloudAppBLB";
import BizDeployNodeConfigFormProviderBaiduCloudBLB from "./BizDeployNodeConfigFormProviderBaiduCloudBLB";
import BizDeployNodeConfigFormProviderBaiduCloudCDN from "./BizDeployNodeConfigFormProviderBaiduCloudCDN";
import BizDeployNodeConfigFormProviderBaishanCDN from "./BizDeployNodeConfigFormProviderBaishanCDN";
import BizDeployNodeConfigFormProviderBaotaPanelConsole from "./BizDeployNodeConfigFormProviderBaotaPanelConsole";
import BizDeployNodeConfigFormProviderBaotaPanelSite from "./BizDeployNodeConfigFormProviderBaotaPanelSite";
import BizDeployNodeConfigFormProviderBaotaWAFSite from "./BizDeployNodeConfigFormProviderBaotaWAFSite";
import BizDeployNodeConfigFormProviderBunnyCDN from "./BizDeployNodeConfigFormProviderBunnyCDN";
import BizDeployNodeConfigFormProviderBytePlusCDN from "./BizDeployNodeConfigFormProviderBytePlusCDN";
import BizDeployNodeConfigFormProviderCdnfly from "./BizDeployNodeConfigFormProviderCdnfly";
import BizDeployNodeConfigFormProviderCTCCCloudAO from "./BizDeployNodeConfigFormProviderCTCCCloudAO";
import BizDeployNodeConfigFormProviderCTCCCloudCDN from "./BizDeployNodeConfigFormProviderCTCCCloudCDN";
import BizDeployNodeConfigFormProviderCTCCCloudELB from "./BizDeployNodeConfigFormProviderCTCCCloudELB";
import BizDeployNodeConfigFormProviderCTCCCloudICDN from "./BizDeployNodeConfigFormProviderCTCCCloudICDN";
import BizDeployNodeConfigFormProviderCTCCCloudLVDN from "./BizDeployNodeConfigFormProviderCTCCCloudLVDN";
import BizDeployNodeConfigFormProviderDogeCloudCDN from "./BizDeployNodeConfigFormProviderDogeCloudCDN";
import BizDeployNodeConfigFormProviderEdgioApplications from "./BizDeployNodeConfigFormProviderEdgioApplications";
import BizDeployNodeConfigFormProviderFlexCDN from "./BizDeployNodeConfigFormProviderFlexCDN";
import BizDeployNodeConfigFormProviderGcoreCDN from "./BizDeployNodeConfigFormProviderGcoreCDN";
import BizDeployNodeConfigFormProviderGoEdge from "./BizDeployNodeConfigFormProviderGoEdge";
import BizDeployNodeConfigFormProviderHuaweiCloudCDN from "./BizDeployNodeConfigFormProviderHuaweiCloudCDN";
import BizDeployNodeConfigFormProviderHuaweiCloudELB from "./BizDeployNodeConfigFormProviderHuaweiCloudELB";
import BizDeployNodeConfigFormProviderHuaweiCloudWAF from "./BizDeployNodeConfigFormProviderHuaweiCloudWAF";
import BizDeployNodeConfigFormProviderJDCloudALB from "./BizDeployNodeConfigFormProviderJDCloudALB";
import BizDeployNodeConfigFormProviderJDCloudCDN from "./BizDeployNodeConfigFormProviderJDCloudCDN";
import BizDeployNodeConfigFormProviderJDCloudLive from "./BizDeployNodeConfigFormProviderJDCloudLive";
import BizDeployNodeConfigFormProviderJDCloudVOD from "./BizDeployNodeConfigFormProviderJDCloudVOD";
import BizDeployNodeConfigFormProviderKong from "./BizDeployNodeConfigFormProviderKong";
import BizDeployNodeConfigFormProviderKubernetesSecret from "./BizDeployNodeConfigFormProviderKubernetesSecret";
import BizDeployNodeConfigFormProviderLeCDN from "./BizDeployNodeConfigFormProviderLeCDN";
import BizDeployNodeConfigFormProviderLocal from "./BizDeployNodeConfigFormProviderLocal";
import BizDeployNodeConfigFormProviderNetlifySite from "./BizDeployNodeConfigFormProviderNetlifySite";
import BizDeployNodeConfigFormProviderProxmoxVE from "./BizDeployNodeConfigFormProviderProxmoxVE";
import BizDeployNodeConfigFormProviderQiniuCDN from "./BizDeployNodeConfigFormProviderQiniuCDN";
import BizDeployNodeConfigFormProviderQiniuKodo from "./BizDeployNodeConfigFormProviderQiniuKodo";
import BizDeployNodeConfigFormProviderQiniuPili from "./BizDeployNodeConfigFormProviderQiniuPili";
import BizDeployNodeConfigFormProviderRainYunRCDN from "./BizDeployNodeConfigFormProviderRainYunRCDN";
import BizDeployNodeConfigFormProviderRatPanelSite from "./BizDeployNodeConfigFormProviderRatPanelSite";
import BizDeployNodeConfigFormProviderSafeLine from "./BizDeployNodeConfigFormProviderSafeLine";
import BizDeployNodeConfigFormProviderSSH from "./BizDeployNodeConfigFormProviderSSH";
import BizDeployNodeConfigFormProviderTencentCloudCDN from "./BizDeployNodeConfigFormProviderTencentCloudCDN";
import BizDeployNodeConfigFormProviderTencentCloudCLB from "./BizDeployNodeConfigFormProviderTencentCloudCLB";
import BizDeployNodeConfigFormProviderTencentCloudCOS from "./BizDeployNodeConfigFormProviderTencentCloudCOS";
import BizDeployNodeConfigFormProviderTencentCloudCSS from "./BizDeployNodeConfigFormProviderTencentCloudCSS";
import BizDeployNodeConfigFormProviderTencentCloudECDN from "./BizDeployNodeConfigFormProviderTencentCloudECDN";
import BizDeployNodeConfigFormProviderTencentCloudEO from "./BizDeployNodeConfigFormProviderTencentCloudEO";
import BizDeployNodeConfigFormProviderTencentCloudGAAP from "./BizDeployNodeConfigFormProviderTencentCloudGAAP";
import BizDeployNodeConfigFormProviderTencentCloudSCF from "./BizDeployNodeConfigFormProviderTencentCloudSCF";
import BizDeployNodeConfigFormProviderTencentCloudSSL from "./BizDeployNodeConfigFormProviderTencentCloudSSL";
import BizDeployNodeConfigFormProviderTencentCloudSSLDeploy from "./BizDeployNodeConfigFormProviderTencentCloudSSLDeploy";
import BizDeployNodeConfigFormProviderTencentCloudSSLUpdate from "./BizDeployNodeConfigFormProviderTencentCloudSSLUpdate";
import BizDeployNodeConfigFormProviderTencentCloudVOD from "./BizDeployNodeConfigFormProviderTencentCloudVOD";
import BizDeployNodeConfigFormProviderTencentCloudWAF from "./BizDeployNodeConfigFormProviderTencentCloudWAF";
import BizDeployNodeConfigFormProviderUCloudUCDN from "./BizDeployNodeConfigFormProviderUCloudUCDN";
import BizDeployNodeConfigFormProviderUCloudUS3 from "./BizDeployNodeConfigFormProviderUCloudUS3";
import BizDeployNodeConfigFormProviderUniCloudWebHost from "./BizDeployNodeConfigFormProviderUniCloudWebHost";
import BizDeployNodeConfigFormProviderUpyunCDN from "./BizDeployNodeConfigFormProviderUpyunCDN";
import BizDeployNodeConfigFormProviderUpyunFile from "./BizDeployNodeConfigFormProviderUpyunFile";
import BizDeployNodeConfigFormProviderVolcEngineALB from "./BizDeployNodeConfigFormProviderVolcEngineALB";
import BizDeployNodeConfigFormProviderVolcEngineCDN from "./BizDeployNodeConfigFormProviderVolcEngineCDN";
import BizDeployNodeConfigFormProviderVolcEngineCertCenter from "./BizDeployNodeConfigFormProviderVolcEngineCertCenter";
import BizDeployNodeConfigFormProviderVolcEngineCLB from "./BizDeployNodeConfigFormProviderVolcEngineCLB";
import BizDeployNodeConfigFormProviderVolcEngineDCDN from "./BizDeployNodeConfigFormProviderVolcEngineDCDN";
import BizDeployNodeConfigFormProviderVolcEngineImageX from "./BizDeployNodeConfigFormProviderVolcEngineImageX";
import BizDeployNodeConfigFormProviderVolcEngineLive from "./BizDeployNodeConfigFormProviderVolcEngineLive";
import BizDeployNodeConfigFormProviderVolcEngineTOS from "./BizDeployNodeConfigFormProviderVolcEngineTOS";
import BizDeployNodeConfigFormProviderWangsuCDN from "./BizDeployNodeConfigFormProviderWangsuCDN";
import BizDeployNodeConfigFormProviderWangsuCDNPro from "./BizDeployNodeConfigFormProviderWangsuCDNPro";
import BizDeployNodeConfigFormProviderWangsuCertificate from "./BizDeployNodeConfigFormProviderWangsuCertificate";
import BizDeployNodeConfigFormProviderWebhook from "./BizDeployNodeConfigFormProviderWebhook";
import { NodeType } from "../nodes/typings";

export interface BizDeployNodeConfigFormProps {
  form: FormInstance;
  node: FlowNodeEntity;
}

const BizDeployNodeConfigForm = ({ node, ...props }: BizDeployNodeConfigFormProps) => {
  if (node.flowNodeType !== NodeType.BizDeploy) {
    console.warn(`[certimate] current workflow node type is not: ${NodeType.BizDeploy}`);
  }

  const { i18n, t } = useTranslation();

  const { token: themeToken } = theme.useToken();

  const { getWorkflowOuptutBeforeId } = useWorkflowStore(useZustandShallowSelector(["updateNode", "getWorkflowOuptutBeforeId"]));

  const initialValues = useMemo(() => {
    return getNodeForm(node)?.getValueIn("config") as WorkflowNodeConfigForDeploy | undefined;
  }, [node]);

  const formSchema = getSchema({ i18n });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm({
    form: props.form,
    name: "workflowNodeBizDeployConfigForm",
    initialValues: initialValues ?? getInitialValues(),
  });

  const fieldProvider = Form.useWatch<string>("provider", { form: formInst, preserve: true });

  const [showProviderAccess, setShowProviderAccess] = useState(false);
  useEffect(() => {
    // 内置的部署提供商（如本地部署）无需显示授权信息字段
    if (fieldProvider) {
      const provider = deploymentProvidersMap.get(fieldProvider);
      setShowProviderAccess(!provider?.builtin);
    } else {
      setShowProviderAccess(false);
    }
  }, [fieldProvider]);

  // TODO: 适配 flowgram
  const certificateCandidates = useMemo(() => {
    const previousNodes = getWorkflowOuptutBeforeId(node.id, "certificate");
    return previousNodes
      .filter((node) => node.type === WorkflowNodeType.Apply || node.type === WorkflowNodeType.Upload)
      .map((item) => {
        return {
          label: item.name,
          options: (item.outputs ?? [])?.map((output) => {
            return {
              label: output.label,
              value: `${item.id}#${output.name}`,
            };
          }),
        };
      })
      .filter((group) => group.options.length > 0);
  }, [node.id]);

  const NestedProviderConfigFields = useMemo(() => {
    /*
        注意：如果追加新的子组件，请保持以 ASCII 排序。
        NOTICE: If you add new child component, please keep ASCII order.
       */
    switch (fieldProvider) {
      case DEPLOYMENT_PROVIDERS["1PANEL_CONSOLE"]:
        return BizDeployNodeConfigFormProvider1PanelConsole;
      case DEPLOYMENT_PROVIDERS["1PANEL_SITE"]:
        return BizDeployNodeConfigFormProvider1PanelSite;
      case DEPLOYMENT_PROVIDERS.ALIYUN_ALB:
        return BizDeployNodeConfigFormProviderAliyunALB;
      case DEPLOYMENT_PROVIDERS.ALIYUN_APIGW:
        return BizDeployNodeConfigFormProviderAliyunAPIGW;
      case DEPLOYMENT_PROVIDERS.ALIYUN_CAS:
        return BizDeployNodeConfigFormProviderAliyunCAS;
      case DEPLOYMENT_PROVIDERS.ALIYUN_CAS_DEPLOY:
        return BizDeployNodeConfigFormProviderAliyunCASDeploy;
      case DEPLOYMENT_PROVIDERS.ALIYUN_CLB:
        return BizDeployNodeConfigFormProviderAliyunCLB;
      case DEPLOYMENT_PROVIDERS.ALIYUN_CDN:
        return BizDeployNodeConfigFormProviderAliyunCDN;
      case DEPLOYMENT_PROVIDERS.ALIYUN_DCDN:
        return BizDeployNodeConfigFormProviderAliyunDCDN;
      case DEPLOYMENT_PROVIDERS.ALIYUN_DDOS:
        return BizDeployNodeConfigFormProviderAliyunDDoS;
      case DEPLOYMENT_PROVIDERS.ALIYUN_ESA:
        return BizDeployNodeConfigFormProviderAliyunESA;
      case DEPLOYMENT_PROVIDERS.ALIYUN_FC:
        return BizDeployNodeConfigFormProviderAliyunFC;
      case DEPLOYMENT_PROVIDERS.ALIYUN_GA:
        return BizDeployNodeConfigFormProviderAliyunGA;
      case DEPLOYMENT_PROVIDERS.ALIYUN_LIVE:
        return BizDeployNodeConfigFormProviderAliyunLive;
      case DEPLOYMENT_PROVIDERS.ALIYUN_NLB:
        return BizDeployNodeConfigFormProviderAliyunNLB;
      case DEPLOYMENT_PROVIDERS.ALIYUN_OSS:
        return BizDeployNodeConfigFormProviderAliyunOSS;
      case DEPLOYMENT_PROVIDERS.ALIYUN_VOD:
        return BizDeployNodeConfigFormProviderAliyunVOD;
      case DEPLOYMENT_PROVIDERS.ALIYUN_WAF:
        return BizDeployNodeConfigFormProviderAliyunWAF;
      case DEPLOYMENT_PROVIDERS.APISIX:
        return BizDeployNodeConfigFormProviderAPISIX;
      case DEPLOYMENT_PROVIDERS.AWS_ACM:
        return BizDeployNodeConfigFormProviderAWSACM;
      case DEPLOYMENT_PROVIDERS.AWS_CLOUDFRONT:
        return BizDeployNodeConfigFormProviderAWSCloudFront;
      case DEPLOYMENT_PROVIDERS.AWS_IAM:
        return BizDeployNodeConfigFormProviderAWSIAM;
      case DEPLOYMENT_PROVIDERS.AZURE_KEYVAULT:
        return BizDeployNodeConfigFormProviderAzureKeyVault;
      case DEPLOYMENT_PROVIDERS.BAIDUCLOUD_APPBLB:
        return BizDeployNodeConfigFormProviderBaiduCloudAppBLB;
      case DEPLOYMENT_PROVIDERS.BAIDUCLOUD_BLB:
        return BizDeployNodeConfigFormProviderBaiduCloudBLB;
      case DEPLOYMENT_PROVIDERS.BAIDUCLOUD_CDN:
        return BizDeployNodeConfigFormProviderBaiduCloudCDN;
      case DEPLOYMENT_PROVIDERS.BAISHAN_CDN:
        return BizDeployNodeConfigFormProviderBaishanCDN;
      case DEPLOYMENT_PROVIDERS.BAOTAPANEL_CONSOLE:
        return BizDeployNodeConfigFormProviderBaotaPanelConsole;
      case DEPLOYMENT_PROVIDERS.BAOTAPANEL_SITE:
        return BizDeployNodeConfigFormProviderBaotaPanelSite;
      case DEPLOYMENT_PROVIDERS.BAOTAWAF_SITE:
        return BizDeployNodeConfigFormProviderBaotaWAFSite;
      case DEPLOYMENT_PROVIDERS.BUNNY_CDN:
        return BizDeployNodeConfigFormProviderBunnyCDN;
      case DEPLOYMENT_PROVIDERS.BYTEPLUS_CDN:
        return BizDeployNodeConfigFormProviderBytePlusCDN;
      case DEPLOYMENT_PROVIDERS.CDNFLY:
        return BizDeployNodeConfigFormProviderCdnfly;
      case DEPLOYMENT_PROVIDERS.CTCCCLOUD_AO:
        return BizDeployNodeConfigFormProviderCTCCCloudAO;
      case DEPLOYMENT_PROVIDERS.CTCCCLOUD_CDN:
        return BizDeployNodeConfigFormProviderCTCCCloudCDN;
      case DEPLOYMENT_PROVIDERS.CTCCCLOUD_ELB:
        return BizDeployNodeConfigFormProviderCTCCCloudELB;
      case DEPLOYMENT_PROVIDERS.CTCCCLOUD_ICDN:
        return BizDeployNodeConfigFormProviderCTCCCloudICDN;
      case DEPLOYMENT_PROVIDERS.CTCCCLOUD_LVDN:
        return BizDeployNodeConfigFormProviderCTCCCloudLVDN;
      case DEPLOYMENT_PROVIDERS.DOGECLOUD_CDN:
        return BizDeployNodeConfigFormProviderDogeCloudCDN;
      case DEPLOYMENT_PROVIDERS.EDGIO_APPLICATIONS:
        return BizDeployNodeConfigFormProviderEdgioApplications;
      case DEPLOYMENT_PROVIDERS.FLEXCDN:
        return BizDeployNodeConfigFormProviderFlexCDN;
      case DEPLOYMENT_PROVIDERS.GCORE_CDN:
        return BizDeployNodeConfigFormProviderGcoreCDN;
      case DEPLOYMENT_PROVIDERS.GOEDGE:
        return BizDeployNodeConfigFormProviderGoEdge;
      case DEPLOYMENT_PROVIDERS.HUAWEICLOUD_CDN:
        return BizDeployNodeConfigFormProviderHuaweiCloudCDN;
      case DEPLOYMENT_PROVIDERS.HUAWEICLOUD_ELB:
        return BizDeployNodeConfigFormProviderHuaweiCloudELB;
      case DEPLOYMENT_PROVIDERS.HUAWEICLOUD_WAF:
        return BizDeployNodeConfigFormProviderHuaweiCloudWAF;
      case DEPLOYMENT_PROVIDERS.JDCLOUD_ALB:
        return BizDeployNodeConfigFormProviderJDCloudALB;
      case DEPLOYMENT_PROVIDERS.JDCLOUD_CDN:
        return BizDeployNodeConfigFormProviderJDCloudCDN;
      case DEPLOYMENT_PROVIDERS.JDCLOUD_LIVE:
        return BizDeployNodeConfigFormProviderJDCloudLive;
      case DEPLOYMENT_PROVIDERS.JDCLOUD_VOD:
        return BizDeployNodeConfigFormProviderJDCloudVOD;
      case DEPLOYMENT_PROVIDERS.KONG:
        return BizDeployNodeConfigFormProviderKong;
      case DEPLOYMENT_PROVIDERS.KUBERNETES_SECRET:
        return BizDeployNodeConfigFormProviderKubernetesSecret;
      case DEPLOYMENT_PROVIDERS.LECDN:
        return BizDeployNodeConfigFormProviderLeCDN;
      case DEPLOYMENT_PROVIDERS.LOCAL:
        return BizDeployNodeConfigFormProviderLocal;
      case DEPLOYMENT_PROVIDERS.NETLIFY_SITE:
        return BizDeployNodeConfigFormProviderNetlifySite;
      case DEPLOYMENT_PROVIDERS.PROXMOXVE:
        return BizDeployNodeConfigFormProviderProxmoxVE;
      case DEPLOYMENT_PROVIDERS.QINIU_CDN:
        return BizDeployNodeConfigFormProviderQiniuCDN;
      case DEPLOYMENT_PROVIDERS.QINIU_KODO:
        return BizDeployNodeConfigFormProviderQiniuKodo;
      case DEPLOYMENT_PROVIDERS.QINIU_PILI:
        return BizDeployNodeConfigFormProviderQiniuPili;
      case DEPLOYMENT_PROVIDERS.RAINYUN_RCDN:
        return BizDeployNodeConfigFormProviderRainYunRCDN;
      case DEPLOYMENT_PROVIDERS.RATPANEL_SITE:
        return BizDeployNodeConfigFormProviderRatPanelSite;
      case DEPLOYMENT_PROVIDERS.SAFELINE:
        return BizDeployNodeConfigFormProviderSafeLine;
      case DEPLOYMENT_PROVIDERS.SSH:
        return BizDeployNodeConfigFormProviderSSH;
      case DEPLOYMENT_PROVIDERS.TENCENTCLOUD_CDN:
        return BizDeployNodeConfigFormProviderTencentCloudCDN;
      case DEPLOYMENT_PROVIDERS.TENCENTCLOUD_CLB:
        return BizDeployNodeConfigFormProviderTencentCloudCLB;
      case DEPLOYMENT_PROVIDERS.TENCENTCLOUD_COS:
        return BizDeployNodeConfigFormProviderTencentCloudCOS;
      case DEPLOYMENT_PROVIDERS.TENCENTCLOUD_CSS:
        return BizDeployNodeConfigFormProviderTencentCloudCSS;
      case DEPLOYMENT_PROVIDERS.TENCENTCLOUD_ECDN:
        return BizDeployNodeConfigFormProviderTencentCloudECDN;
      case DEPLOYMENT_PROVIDERS.TENCENTCLOUD_EO:
        return BizDeployNodeConfigFormProviderTencentCloudEO;
      case DEPLOYMENT_PROVIDERS.TENCENTCLOUD_GAAP:
        return BizDeployNodeConfigFormProviderTencentCloudGAAP;
      case DEPLOYMENT_PROVIDERS.TENCENTCLOUD_SCF:
        return BizDeployNodeConfigFormProviderTencentCloudSCF;
      case DEPLOYMENT_PROVIDERS.TENCENTCLOUD_SSL:
        return BizDeployNodeConfigFormProviderTencentCloudSSL;
      case DEPLOYMENT_PROVIDERS.TENCENTCLOUD_SSL_DEPLOY:
        return BizDeployNodeConfigFormProviderTencentCloudSSLDeploy;
      case DEPLOYMENT_PROVIDERS.TENCENTCLOUD_SSL_UPDATE:
        return BizDeployNodeConfigFormProviderTencentCloudSSLUpdate;
      case DEPLOYMENT_PROVIDERS.TENCENTCLOUD_VOD:
        return BizDeployNodeConfigFormProviderTencentCloudVOD;
      case DEPLOYMENT_PROVIDERS.TENCENTCLOUD_WAF:
        return BizDeployNodeConfigFormProviderTencentCloudWAF;
      case DEPLOYMENT_PROVIDERS.UCLOUD_UCDN:
        return BizDeployNodeConfigFormProviderUCloudUCDN;
      case DEPLOYMENT_PROVIDERS.UCLOUD_US3:
        return BizDeployNodeConfigFormProviderUCloudUS3;
      case DEPLOYMENT_PROVIDERS.UNICLOUD_WEBHOST:
        return BizDeployNodeConfigFormProviderUniCloudWebHost;
      case DEPLOYMENT_PROVIDERS.UPYUN_CDN:
        return BizDeployNodeConfigFormProviderUpyunCDN;
      case DEPLOYMENT_PROVIDERS.UPYUN_FILE:
        return BizDeployNodeConfigFormProviderUpyunFile;
      case DEPLOYMENT_PROVIDERS.VOLCENGINE_ALB:
        return BizDeployNodeConfigFormProviderVolcEngineALB;
      case DEPLOYMENT_PROVIDERS.VOLCENGINE_CDN:
        return BizDeployNodeConfigFormProviderVolcEngineCDN;
      case DEPLOYMENT_PROVIDERS.VOLCENGINE_CERTCENTER:
        return BizDeployNodeConfigFormProviderVolcEngineCertCenter;
      case DEPLOYMENT_PROVIDERS.VOLCENGINE_CLB:
        return BizDeployNodeConfigFormProviderVolcEngineCLB;
      case DEPLOYMENT_PROVIDERS.VOLCENGINE_DCDN:
        return BizDeployNodeConfigFormProviderVolcEngineDCDN;
      case DEPLOYMENT_PROVIDERS.VOLCENGINE_IMAGEX:
        return BizDeployNodeConfigFormProviderVolcEngineImageX;
      case DEPLOYMENT_PROVIDERS.VOLCENGINE_LIVE:
        return BizDeployNodeConfigFormProviderVolcEngineLive;
      case DEPLOYMENT_PROVIDERS.VOLCENGINE_TOS:
        return BizDeployNodeConfigFormProviderVolcEngineTOS;
      case DEPLOYMENT_PROVIDERS.WANGSU_CDN:
        return BizDeployNodeConfigFormProviderWangsuCDN;
      case DEPLOYMENT_PROVIDERS.WANGSU_CDNPRO:
        return BizDeployNodeConfigFormProviderWangsuCDNPro;
      case DEPLOYMENT_PROVIDERS.WANGSU_CERTIFICATE:
        return BizDeployNodeConfigFormProviderWangsuCertificate;
      case DEPLOYMENT_PROVIDERS.WEBHOOK:
        return BizDeployNodeConfigFormProviderWebhook;
    }
  }, [fieldProvider]);

  const handleProviderPick = (value: string) => {
    formInst.setFieldValue("provider", value);
    console.log();
  };

  const handleProviderSelect = (value?: string | undefined) => {
    // 切换部署目标时重置表单，避免其他部署目标的配置字段影响当前部署目标
    if (initialValues?.provider === value) {
      formInst.resetFields();
    } else {
      const oldValues = formInst.getFieldsValue();
      const newValues: Record<string, unknown> = {};
      for (const key in oldValues) {
        if (key === "provider" || key === "providerAccessId" || key === "certificate" || key === "skipOnLastSucceeded") {
          newValues[key] = oldValues[key];
        } else {
          delete newValues[key];
        }
      }
      formInst.setFieldsValue(newValues);

      if (deploymentProvidersMap.get(fieldProvider)?.provider !== deploymentProvidersMap.get(value!)?.provider) {
        formInst.setFieldValue("providerAccessId", void 0);
      }
    }
  };

  return (
    <NodeFormContextProvider value={{ node }}>
      <Form {...formProps} clearOnDestroy={true} form={formInst} layout="vertical" preserve={false} scrollToFirstError>
        <Show when={!fieldProvider}>
          <DeploymentProviderPicker autoFocus placeholder={t("workflow_node.deploy.form.provider.search.placeholder")} onSelect={handleProviderPick} />
        </Show>

        <div style={{ display: fieldProvider ? "block" : "none" }}>
          <div id="parameters" data-anchor="parameters">
            <Form.Item name="provider" label={t("workflow_node.deploy.form.provider.label")} rules={[formRule]}>
              <DeploymentProviderSelect
                allowClear
                disabled={!!initialValues?.provider}
                placeholder={t("workflow_node.deploy.form.provider.placeholder")}
                showSearch
                onSelect={handleProviderSelect}
                onClear={handleProviderSelect}
              />
            </Form.Item>

            <Form.Item
              className="relative"
              hidden={!showProviderAccess}
              label={t("workflow_node.deploy.form.provider_access.label")}
              tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.provider_access.tooltip") }}></span>}
            >
              <div className="absolute -top-[6px] right-0 -translate-y-full">
                <AccessEditDrawer
                  data={{ provider: deploymentProvidersMap.get(fieldProvider!)?.provider }}
                  mode="create"
                  trigger={
                    <Button size="small" type="link">
                      {t("workflow_node.deploy.form.provider_access.button")}
                      <IconPlus size="1.25em" />
                    </Button>
                  }
                  usage="hosting"
                  afterSubmit={(record) => {
                    const provider = accessProvidersMap.get(record.provider);
                    if (provider?.usages?.includes(ACCESS_USAGES.HOSTING)) {
                      formInst.setFieldValue("providerAccessId", record.id);
                    }
                  }}
                />
              </div>
              <Form.Item name="providerAccessId" rules={[formRule]} noStyle>
                <AccessSelect
                  placeholder={t("workflow_node.deploy.form.provider_access.placeholder")}
                  showSearch
                  onFilter={(_, option) => {
                    if (option.reserve) return false;
                    if (fieldProvider) return deploymentProvidersMap.get(fieldProvider)?.provider === option.provider;

                    const provider = accessProvidersMap.get(option.provider);
                    return !!provider?.usages?.includes(ACCESS_USAGES.HOSTING);
                  }}
                />
              </Form.Item>
            </Form.Item>

            <Form.Item
              name="certificate"
              label={t("workflow_node.deploy.form.certificate.label")}
              rules={[formRule]}
              tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.certificate.tooltip") }}></span>}
            >
              <Select
                labelRender={({ label, value }) => {
                  if (value != null) {
                    const group = certificateCandidates.find((group) => group.options.some((option) => option.value === value));
                    return `${group?.label} - ${label}`;
                  }

                  return <span style={{ color: themeToken.colorTextPlaceholder }}>{t("workflow_node.deploy.form.certificate.placeholder")}</span>;
                }}
                options={certificateCandidates}
                placeholder={t("workflow_node.deploy.form.certificate.placeholder")}
              />
            </Form.Item>
          </div>

          <div id="deployment" data-anchor="deployment">
            <FormNestedFieldsContextProvider value={{ parentNamePath: "providerConfig" }}>
              {NestedProviderConfigFields && (
                <>
                  <Divider size="small">
                    <Typography.Text className="text-xs font-normal" type="secondary">
                      {t("workflow_node.deploy.form_anchor.deployment.title")}
                    </Typography.Text>
                  </Divider>

                  <NestedProviderConfigFields />
                </>
              )}
            </FormNestedFieldsContextProvider>
          </div>

          <div id="strategy" data-anchor="strategy">
            <Divider size="small">
              <Typography.Text className="text-xs font-normal" type="secondary">
                {t("workflow_node.deploy.form_anchor.strategy.title")}
              </Typography.Text>
            </Divider>

            <Form.Item label={t("workflow_node.deploy.form.skip_on_last_succeeded.label")}>
              <Flex align="center" gap={8} wrap="wrap">
                <div>{t("workflow_node.deploy.form.skip_on_last_succeeded.prefix")}</div>
                <Form.Item name="skipOnLastSucceeded" noStyle rules={[formRule]}>
                  <Switch
                    checkedChildren={t("workflow_node.deploy.form.skip_on_last_succeeded.switch.on")}
                    unCheckedChildren={t("workflow_node.deploy.form.skip_on_last_succeeded.switch.off")}
                  />
                </Form.Item>
                <div>{t("workflow_node.deploy.form.skip_on_last_succeeded.suffix")}</div>
              </Flex>
            </Form.Item>
          </div>
        </div>
      </Form>
    </NodeFormContextProvider>
  );
};

const getAnchorItems = ({ i18n = getI18n() }: { i18n?: ReturnType<typeof getI18n> }): Required<AnchorProps>["items"] => {
  const { t } = i18n;

  return ["parameters", "deployment", "strategy"].map((key) => ({
    key: key,
    title: t(`workflow_node.deploy.form_anchor.${key}.tab`),
    href: "#" + key,
  }));
};

const getInitialValues = (): Nullish<z.infer<ReturnType<typeof getSchema>>> => {
  return {
    certificate: "",
    provider: "",
    ...defaultNodeConfigForDeploy(),
  };
};

const getSchema = ({ i18n = getI18n() }: { i18n?: ReturnType<typeof getI18n> }) => {
  const { t } = i18n;

  return z
    .object({
      certificate: z.string().nonempty(t("workflow_node.deploy.form.certificate.placeholder")),
      provider: z.string().nonempty(t("workflow_node.deploy.form.provider.placeholder")),
      providerAccessId: z.string().nullish(),
      providerConfig: z.any().nullish(),
      skipOnLastSucceeded: z.boolean().nullish(),
    })
    .superRefine((values, ctx) => {
      if (values.provider) {
        const provider = deploymentProvidersMap.get(values.provider);
        if (!provider?.builtin && !values.providerAccessId) {
          ctx.addIssue({
            code: "custom",
            message: t("workflow_node.deploy.form.provider_access.placeholder"),
            path: ["providerAccessId"],
          });
        }
      }
    });
};

const _default = Object.assign(BizDeployNodeConfigForm, {
  getAnchorItems,
  getSchema,
});

export default _default;
