#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
import sys
from datetime import datetime

# 设置输出编码为 UTF-8
sys.stdout.reconfigure(encoding='utf-8')

# 前端文件功能描述映射
FRONTEND_DESCRIPTIONS = {
    # Admin 目录
    'admin/src/main.ts': 'Admin 应用入口文件',
    'admin/src/App.vue': 'Admin 根组件',
    'admin/src/router/index.ts': '路由配置',
    'admin/src/store/index.ts': '状态管理',
    'admin/src/api': 'API 接口定义',
    'admin/src/views': '页面组件',
    'admin/src/components': '公共组件',
    'admin/src/utils': '工具函数',
    'admin/src/types': '类型定义',
    'admin/src/composables': '组合式函数',
    'admin/src/directives': '自定义指令',
    'admin/src/plugins': '插件配置',
    'admin/src/styles': '样式文件',
    'admin/src/assets': '静态资源',
    'admin/vite.config.ts': 'Vite 构建配置',
    'admin/tsconfig.json': 'TypeScript 配置',
    'admin/package.json': '项目依赖配置',

    # Blog 目录
    'blog/src/main.ts': 'Blog 应用入口文件',
    'blog/src/App.vue': 'Blog 根组件',
    'blog/src/router/index.ts': '路由配置',
    'blog/src/store/index.ts': '状态管理',
    'blog/src/api': 'API 接口定义',
    'blog/src/views': '页面组件',
    'blog/src/components': '公共组件',
    'blog/src/utils': '工具函数',
    'blog/src/types': '类型定义',
    'blog/src/composables': '组合式函数',
    'blog/src/directives': '自定义指令',
    'blog/src/plugins': '插件配置',
    'blog/src/styles': '样式文件',
    'blog/src/assets': '静态资源',
    'blog/vite.config.ts': 'Vite 构建配置',
    'blog/tsconfig.json': 'TypeScript 配置',
    'blog/package.json': '项目依赖配置',
}

def get_description(file_path):
    """根据文件路径获取功能描述"""
    # 标准化路径
    normalized_path = file_path.replace('\\', '/')

    # 精确匹配
    if normalized_path in FRONTEND_DESCRIPTIONS:
        return FRONTEND_DESCRIPTIONS[normalized_path]

    # 目录匹配
    for pattern, desc in FRONTEND_DESCRIPTIONS.items():
        if '/' in pattern and not pattern.endswith('.ts') and not pattern.endswith('.vue') and not pattern.endswith('.json'):
            if pattern in normalized_path:
                # 根据文件名细化描述
                filename = os.path.basename(file_path)
                if 'api' in pattern:
                    return f'{desc} - {filename.replace(".ts", "")}'
                elif 'views' in pattern:
                    return f'{desc} - {filename.replace(".vue", "")}页面'
                elif 'components' in pattern:
                    return f'{desc} - {filename.replace(".vue", "")}组件'
                elif 'utils' in pattern:
                    return f'{desc} - {filename.replace(".ts", "")}工具'
                elif 'types' in pattern:
                    return f'{desc} - {filename.replace(".ts", "")}类型'
                elif 'composables' in pattern:
                    return f'{desc} - {filename.replace(".ts", "")}组合函数'
                return desc

    # 默认描述
    if file_path.endswith('.vue'):
        return 'Vue 组件'
    elif file_path.endswith('.ts'):
        return 'TypeScript 模块'
    elif file_path.endswith('.js'):
        return 'JavaScript 模块'
    elif file_path.endswith('.css') or file_path.endswith('.scss'):
        return '样式文件'
    else:
        return '前端文件'

def has_header_comment(content, file_ext):
    """检查文件是否已有头部注释"""
    lines = content.strip().split('\n')
    if not lines:
        return False

    first_line = lines[0].strip()

    if file_ext in ['.vue']:
        # Vue 文件检查 <!-- 注释
        return first_line.startswith('<!--') and '项目名称' in content[:500]
    else:
        # TS/JS 文件检查 /* 注释
        return first_line.startswith('/*') and '项目名称' in content[:500]

def create_header_comment(file_path, file_ext):
    """创建头部注释"""
    filename = os.path.basename(file_path)
    current_time = datetime.now().strftime('%Y-%m-%d %H:%M:%S')
    description = get_description(file_path)

    if file_ext == '.vue':
        # Vue 文件使用 HTML 注释
        return f"""<!--
项目名称：JeriBlog
文件名称：{filename}
创建时间：{current_time}

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：{description}
-->

"""
    else:
        # TS/JS 文件使用 /* */ 注释
        return f"""/*
项目名称：JeriBlog
文件名称：{filename}
创建时间：{current_time}

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：{description}
*/

"""

def process_file(file_path):
    """处理单个文件"""
    try:
        file_ext = os.path.splitext(file_path)[1].lower()

        # 只处理 .vue, .ts, .js 文件
        if file_ext not in ['.vue', '.ts', '.js']:
            return False

        # 读取文件内容
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()

        # 检查是否已有头部注释
        if has_header_comment(content, file_ext):
            print(f"⏭️  跳过（已有注释）: {file_path}")
            return False

        # 创建头部注释
        header = create_header_comment(file_path, file_ext)

        # 写入文件
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(header + content)

        print(f"✅ 已处理: {file_path}")
        return True

    except Exception as e:
        print(f"❌ 处理失败 {file_path}: {str(e)}")
        return False

def process_directory(directory):
    """处理目录下的所有文件"""
    processed_count = 0
    skipped_count = 0

    for root, dirs, files in os.walk(directory):
        # 跳过 node_modules 和 dist 目录
        dirs[:] = [d for d in dirs if d not in ['node_modules', 'dist', '.git']]

        for file in files:
            file_path = os.path.join(root, file)
            if process_file(file_path):
                processed_count += 1
            else:
                if os.path.splitext(file)[1].lower() in ['.vue', '.ts', '.js']:
                    skipped_count += 1

    return processed_count, skipped_count

def main():
    if len(sys.argv) < 2:
        print("用法: python add_frontend_headers.py <目录路径>")
        sys.exit(1)

    directory = sys.argv[1]

    if not os.path.exists(directory):
        print(f"❌ 目录不存在: {directory}")
        sys.exit(1)

    print(f"开始处理目录: {directory}")
    print("=" * 60)

    processed, skipped = process_directory(directory)

    print("=" * 60)
    print(f"✅ 处理完成！")
    print(f"   新增注释: {processed} 个文件")
    print(f"   跳过文件: {skipped} 个文件")

if __name__ == '__main__':
    main()
