#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
项目名称：JeriBlog
文件名称：convert_favicon.py
创建时间：2026-04-16 16:00:00

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：将 favicon.ico 转换为不同尺寸的 PNG 图标
"""

from PIL import Image
import os

def convert_favicon_to_png(ico_path, output_dir):
    """
    将 favicon.ico 转换为不同尺寸的 PNG 文件

    Args:
        ico_path: favicon.ico 文件路径
        output_dir: 输出目录
    """
    if not os.path.exists(ico_path):
        print(f"[ERROR] 找不到文件: {ico_path}")
        return

    try:
        # 打开 ICO 文件
        img = Image.open(ico_path)

        # 如果 ICO 包含多个尺寸，选择最大的
        if hasattr(img, 'size'):
            print(f"[INFO] 原始尺寸: {img.size}")

        # 转换为 RGBA 模式（支持透明度）
        if img.mode != 'RGBA':
            img = img.convert('RGBA')

        # 生成 apple-touch-icon.png (180x180)
        apple_icon = img.resize((180, 180), Image.Resampling.LANCZOS)
        apple_icon_path = os.path.join(output_dir, 'apple-touch-icon.png')
        apple_icon.save(apple_icon_path, 'PNG', optimize=True)
        print(f"[OK] 已生成: {apple_icon_path} (180x180)")

        # 生成 pwa-192x192.png (192x192)
        pwa_icon = img.resize((192, 192), Image.Resampling.LANCZOS)
        pwa_icon_path = os.path.join(output_dir, 'pwa-192x192.png')
        pwa_icon.save(pwa_icon_path, 'PNG', optimize=True)
        print(f"[OK] 已生成: {pwa_icon_path} (192x192)")

        # 生成 pwa-512x512.png (512x512) - PWA 标准尺寸
        pwa_large_icon = img.resize((512, 512), Image.Resampling.LANCZOS)
        pwa_large_icon_path = os.path.join(output_dir, 'pwa-512x512.png')
        pwa_large_icon.save(pwa_large_icon_path, 'PNG', optimize=True)
        print(f"[OK] 已生成: {pwa_large_icon_path} (512x512)")

        print("\n[SUCCESS] 转换完成!")

    except Exception as e:
        print(f"[ERROR] 转换失败: {str(e)}")

if __name__ == '__main__':
    # admin 目录
    admin_ico = r'D:\GitHub\我的项目\JeriBlog\admin\public\favicon.ico'
    admin_output = r'D:\GitHub\我的项目\JeriBlog\admin\public'

    print("=" * 60)
    print("开始转换 Admin 图标...")
    print("=" * 60)
    convert_favicon_to_png(admin_ico, admin_output)

    print("\n" + "=" * 60)
    print("开始转换 Blog 图标...")
    print("=" * 60)

    # blog 目录
    blog_ico = r'D:\GitHub\我的项目\JeriBlog\blog\public\favicon.ico'
    blog_output = r'D:\GitHub\我的项目\JeriBlog\blog\public'
    convert_favicon_to_png(blog_ico, blog_output)
