const TARGET_ATTRIBUTES = new Set(["href", "src"]);

const getStaticPrefix = (node) => {
  if (!node) {
    return;
  }

  if (node.type === "Literal") {
    return typeof node.value === "string" ? node.value : void 0;
  }

  if (node.type === "TemplateLiteral") {
    return node.quasis[0]?.value.raw;
  }

  if (node.type === "BinaryExpression" && node.operator === "+") {
    return getStaticPrefix(node.left);
  }
};

export default {
  meta: {
    type: "problem",
    docs: {
      description: "Disallow root-relative URLs in JSX href and src attributes",
    },
    messages: {
      useBasePath: "Root-relative {{attribute}} bypasses the application base path. Wrap the URL with withBasePath().",
    },
    schema: [],
  },
  create(context) {
    return {
      JSXAttribute(node) {
        const attribute = node.name.type === "JSXIdentifier" ? node.name.name : "";
        if (!TARGET_ATTRIBUTES.has(attribute)) {
          return;
        }

        const value = node.value?.type === "JSXExpressionContainer" ? getStaticPrefix(node.value.expression) : getStaticPrefix(node.value);
        if (value?.startsWith("/") && !value.startsWith("//")) {
          context.report({
            node: node.value ?? node,
            messageId: "useBasePath",
            data: { attribute },
          });
        }
      },
    };
  },
};
