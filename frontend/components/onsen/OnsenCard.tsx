import React from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { FaStar, FaMapMarkerAlt, FaCalendarAlt, FaWater } from 'react-icons/fa';
import { OnsenLog } from '@/types';

interface OnsenCardProps {
  onsenLog: OnsenLog;
}

const OnsenCard: React.FC<OnsenCardProps> = ({ onsenLog }) => {
  // 日付をフォーマット
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('ja-JP', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  // 評価を星で表示
  const renderRating = (rating: number | undefined) => {
    if (rating === undefined) return null;
    
    return (
      <div className="flex items-center">
        {[...Array(5)].map((_, i) => (
          <FaStar
            key={i}
            className={`w-4 h-4 ${
              i < rating ? 'text-yellow-400' : 'text-gray-300'
            }`}
          />
        ))}
        <span className="ml-1 text-sm text-gray-600">{rating}</span>
      </div>
    );
  };

  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow duration-300">
      <Link href={`/onsen/${onsenLog.id}`}>
        <div className="p-4">
          <div className="flex justify-between items-start mb-2">
            <h3 className="text-xl font-semibold text-gray-800 hover:text-onsen transition-colors duration-200">
              {onsenLog.name}
            </h3>
            {renderRating(onsenLog.rating)}
          </div>
          
          {onsenLog.location && (
            <div className="flex items-center text-gray-600 mb-2">
              <FaMapMarkerAlt className="mr-2 text-onsen" />
              <span>{onsenLog.location}</span>
            </div>
          )}
          
          <div className="flex items-center text-gray-600 mb-2">
            <FaCalendarAlt className="mr-2 text-onsen" />
            <span>{formatDate(onsenLog.visit_date)}</span>
          </div>
          
          {onsenLog.spring_type && (
            <div className="flex items-center text-gray-600 mb-2">
              <FaWater className="mr-2 text-onsen" />
              <span>{onsenLog.spring_type}</span>
            </div>
          )}
          
          {onsenLog.features && onsenLog.features.length > 0 && (
            <div className="flex flex-wrap gap-1 mt-3">
              {onsenLog.features.map((feature, index) => (
                <span
                  key={index}
                  className="inline-block bg-gray-100 rounded-full px-3 py-1 text-xs font-semibold text-gray-700"
                >
                  {feature}
                </span>
              ))}
            </div>
          )}
          
          {onsenLog.comment && (
            <div className="mt-3 text-gray-600 line-clamp-2">
              {onsenLog.comment}
            </div>
          )}
        </div>
      </Link>
    </div>
  );
};

export default OnsenCard; 