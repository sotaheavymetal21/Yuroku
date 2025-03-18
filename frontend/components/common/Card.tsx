import React, { ReactNode } from 'react';

interface CardProps {
  children: ReactNode;
  title?: string;
  subtitle?: string;
  footer?: ReactNode;
  className?: string;
  headerClassName?: string;
  bodyClassName?: string;
  footerClassName?: string;
  ariaLabel?: string;
  fullWidth?: boolean;
}

const Card: React.FC<CardProps> = ({
  children,
  title,
  subtitle,
  footer,
  className = '',
  headerClassName = '',
  bodyClassName = '',
  footerClassName = '',
  ariaLabel,
  fullWidth = false,
}) => {
  return (
    <div 
      className={`bg-white rounded-lg shadow-md overflow-hidden ${fullWidth ? 'w-full' : ''} ${className}`}
      aria-label={ariaLabel}
      role={ariaLabel ? 'region' : undefined}
    >
      {(title || subtitle) && (
        <div className={`px-4 sm:px-6 py-4 border-b border-gray-200 ${headerClassName}`}>
          {title && (
            <h3 className="text-lg font-medium text-gray-900 break-words">
              {title}
            </h3>
          )}
          {subtitle && (
            <p className="mt-1 text-sm text-gray-500 break-words">
              {subtitle}
            </p>
          )}
        </div>
      )}
      <div className={`px-4 sm:px-6 py-4 ${bodyClassName}`}>
        {children}
      </div>
      {footer && (
        <div className={`px-4 sm:px-6 py-4 border-t border-gray-200 bg-gray-50 ${footerClassName}`}>
          {footer}
        </div>
      )}
    </div>
  );
};

export default Card; 